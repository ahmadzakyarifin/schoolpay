package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(appEnv, frontendURL string) gin.HandlerFunc {
	return cors.New(cors.Config{
		// perintah apa saja yang boleh dilakukan backend oleh frontend
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// Header apa saja yang boleh dikirim frontend ke backend.
		AllowHeaders: []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "X-Idempotency-Key", "X-Device-ID", "X-App-Platform", "X-App-Version"},
		// Header response apa saja yang boleh dibaca oleh JavaScript frontend.
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type"},
		// browser boleh membawa credential saat request lintas origin.
		AllowCredentials: true,
		// Browser boleh menyimpan hasil pengecekan OPTIONS selama 12 jam.
		MaxAge: 12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == frontendURL
		},
	})
}
