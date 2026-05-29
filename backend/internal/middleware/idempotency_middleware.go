package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type idempotencyResponse struct {
	StatusCode  int         `json:"status_code"`
	ContentType string      `json:"content_type"`
	Body        json.RawMessage `json:"body"`
}

type responseCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// IdempotencyMiddleware caches successful unsafe responses for requests carrying
// X-Idempotency-Key. It is used by the PWA offline outbox to safely retry sync.
func IdempotencyMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		key := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
		if key == "" || rdb == nil {
			c.Next()
			return
		}

		ctx := context.Background()
		redisKey := "idempotency:" + key

		if cached, err := rdb.Get(ctx, redisKey).Result(); err == nil && cached != "" {
			var res idempotencyResponse
			if json.Unmarshal([]byte(cached), &res) == nil {
				contentType := res.ContentType
				if contentType == "" {
					contentType = "application/json; charset=utf-8"
				}
				c.Data(res.StatusCode, contentType, res.Body)
				c.Abort()
				return
			}
		}

		capture := &responseCaptureWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = capture
		c.Next()

		status := c.Writer.Status()
		if status < http.StatusOK || status >= http.StatusMultipleChoices || capture.body.Len() == 0 {
			return
		}

		res := idempotencyResponse{
			StatusCode:  status,
			ContentType: c.Writer.Header().Get("Content-Type"),
			Body:        json.RawMessage(capture.body.Bytes()),
		}
		payload, err := json.Marshal(res)
		if err != nil {
			return
		}
		_ = rdb.Set(ctx, redisKey, payload, 24*time.Hour).Err()
	}
}
