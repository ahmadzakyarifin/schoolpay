package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

const maxIdempotencyKeyLen = 255

type idempotencyResponse struct {
	StatusCode  int    `json:"status_code"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type responseCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseCaptureWriter) WriteString(data string) (int, error) {
	w.body.WriteString(data)
	return w.ResponseWriter.WriteString(data)
}

// IdempotencyMiddleware caches successful unsafe responses for requests carrying
// X-Idempotency-Key in the MariaDB database.
func IdempotencyMiddleware(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		key := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
		if key == "" || db == nil {
			c.Next()
			return
		}
		if !isValidIdempotencyKey(key) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "X-Idempotency-Key tidak valid",
			})
			return
		}

		ctx := c.Request.Context()
		requestHash, err := idempotencyRequestHash(c.Request)
		if err != nil {
			log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: failed to read request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "gagal membaca request body untuk idempotency",
			})
			return
		}

		claimed := &notificationdomain.IdempotencyKey{
			Key:         key,
			Status:      notificationdomain.IdempotencyStatusProcessing,
			RequestHash: requestHash,
		}
		if _, err := db.NewInsert().Model(claimed).Exec(ctx); err != nil {
			if !isDuplicateKeyError(err) {
				log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: failed to claim key")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "gagal memproses idempotency key",
				})
				return
			}

			var ik notificationdomain.IdempotencyKey
			err := db.NewSelect().Model(&ik).Where("`key` = ?", key).Scan(ctx)
			if errors.Is(err, sql.ErrNoRows) {
				log.Warn().Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: duplicate key vanished before select")
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "request dengan idempotency key ini sedang diproses",
				})
				return
			}
			if err != nil {
				log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: failed to read existing key")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"message": "gagal membaca idempotency key",
				})
				return
			}

			if ik.RequestHash != "" && ik.RequestHash != requestHash {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "idempotency key sudah digunakan untuk request yang berbeda",
				})
				return
			}

			switch ik.Status {
			case notificationdomain.IdempotencyStatusCompleted:
				if ik.ResponsePayload != "" {
					var res idempotencyResponse
					if json.Unmarshal([]byte(ik.ResponsePayload), &res) == nil {
						contentType := res.ContentType
						if contentType == "" {
							contentType = "application/json; charset=utf-8"
						}
						c.Data(res.StatusCode, contentType, []byte(res.Body))
						c.Abort()
						return
					}
				}

				log.Warn().Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: completed response payload is invalid")
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "response idempotency tidak valid, silakan gunakan key baru",
				})
				return
			case notificationdomain.IdempotencyStatusProcessing:
				c.Header("Retry-After", "2")
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "request dengan idempotency key ini sedang diproses",
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "status idempotency key tidak valid",
				})
				return
			}
		}

		markAsRetryable := func(reason string) {
			if _, err := db.NewDelete().
				Model((*notificationdomain.IdempotencyKey)(nil)).
				Where("`key` = ?", key).
				Where("status = ?", notificationdomain.IdempotencyStatusProcessing).
				Exec(context.WithoutCancel(ctx)); err != nil {
				log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Str("reason", reason).Msg("idempotency: failed to delete retryable key")
			}
		}

		defer func() {
			if rec := recover(); rec != nil {
				markAsRetryable("panic")
				panic(rec)
			}
		}()

		capture := &responseCaptureWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = capture
		c.Next()

		status := c.Writer.Status()
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			markAsRetryable("non_2xx_response")
			return
		}

		res := idempotencyResponse{
			StatusCode:  status,
			ContentType: c.Writer.Header().Get("Content-Type"),
			Body:        capture.body.String(),
		}
		payload, err := json.Marshal(res)
		if err != nil {
			log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: failed to marshal cached response")
			markAsRetryable("marshal_response_failed")
			return
		}

		if _, err := db.NewUpdate().
			Model((*notificationdomain.IdempotencyKey)(nil)).
			Set("status = ?", notificationdomain.IdempotencyStatusCompleted).
			Set("response_payload = ?", string(payload)).
			Set("updated_at = ?", time.Now()).
			Where("`key` = ?", key).
			Where("status = ?", notificationdomain.IdempotencyStatusProcessing).
			Exec(context.WithoutCancel(ctx)); err != nil {
			log.Error().Err(err).Str("key_fp", idempotencyFingerprint(key)).Msg("idempotency: failed to mark key completed")
			markAsRetryable("complete_update_failed")
		}
	}
}

func isValidIdempotencyKey(key string) bool {
	if key == "" || len(key) > maxIdempotencyKeyLen || !utf8.ValidString(key) {
		return false
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_' || r == '.' || r == ':' || r == '~':
		default:
			return false
		}
	}
	return true
}

func idempotencyRequestHash(r *http.Request) (string, error) {
	var body []byte
	if r.Body != nil {
		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s\n%s\n%s\n%s", r.Method, r.URL.Path, r.URL.RawQuery, string(body))))
	return fmt.Sprintf("%x", sum), nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

func idempotencyFingerprint(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:8])
}

func CleanupIdempotencyKeys(ctx context.Context, db *bun.DB, completedTTL time.Duration, processingTTL time.Duration) (int64, error) {
	if db == nil {
		return 0, nil
	}
	if completedTTL <= 0 {
		completedTTL = 24 * time.Hour
	}
	if processingTTL <= 0 {
		processingTTL = 15 * time.Minute
	}

	now := time.Now()
	res, err := db.NewDelete().
		Model((*notificationdomain.IdempotencyKey)(nil)).
		WhereGroup(" AND ", func(q *bun.DeleteQuery) *bun.DeleteQuery {
			return q.
				WhereOr("status = ? AND updated_at < ?", notificationdomain.IdempotencyStatusCompleted, now.Add(-completedTTL)).
				WhereOr("status = ? AND updated_at < ?", notificationdomain.IdempotencyStatusProcessing, now.Add(-processingTTL))
		}).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("idempotency: cleanup failed")
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("idempotency: failed to read cleanup rows affected")
		return 0, err
	}
	if rows > 0 {
		log.Info().Int64("rows", rows).Msg("idempotency: cleaned expired keys")
	}
	return rows, nil
}

func StartIdempotencyCleanupJob(ctx context.Context, db *bun.DB, interval time.Duration, completedTTL time.Duration, processingTTL time.Duration) {
	if db == nil {
		return
	}
	if interval <= 0 {
		interval = time.Hour
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = CleanupIdempotencyKeys(ctx, db, completedTTL, processingTTL)
			}
		}
	}()
}
