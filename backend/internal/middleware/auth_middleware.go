package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string, userRepo repository.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// fmt.Printf("[DEBUG] AuthMiddleware: Method=%s, Path=%s\n", c.Request.Method, c.Request.URL.Path)
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		var tokenStr string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			} else {
				utils.ErrorResponse(c, http.StatusUnauthorized, "format token salah")
				c.Abort()
				return
			}
		} else {
			// Fallback ke query parameter (biasa digunakan untuk WebSocket)
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "token tidak ditemukan")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "token tidak valid atau kadaluarsa")
			c.Abort()
			return
		}

		// Validasi Real-time ke Database
		user, err := userRepo.FindByID(c.Request.Context(), claims.UserID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "pengguna tidak ditemukan atau sudah dihapus")
			c.Abort()
			return
		}

		if !user.IsActive {
			utils.ErrorResponse(c, http.StatusForbidden, "akun anda sedang dinonaktifkan")
			c.Abort()
			return
		}

		// Simpan data user TERBARU dari DB ke context
		c.Set("user_id", user.ID)
		c.Set("role", user.Role)
		c.Set("user_name", user.Name)
		c.Set("user", user)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "role", user.Role)
		ctx = context.WithValue(ctx, "user_name", user.Name)
		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
		ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		userRole, exists := c.Get("role")
		if !exists {
			// fmt.Printf("[CRITICAL] RoleMiddleware: No role in context. Aborting.\n")
			utils.ErrorResponse(c, http.StatusForbidden, "akses ditolak: role tidak ditemukan")
			c.Abort()
			return
		}

		// Convert userRole to string safely
		roleStr := fmt.Sprintf("%v", userRole)
		// fmt.Printf("[DEBUG] RoleMiddleware: userRole=%s, requiredRoles=%v\n", roleStr, roles)

		roleMatch := false
		for _, r := range roles {
			if strings.EqualFold(strings.TrimSpace(roleStr), strings.TrimSpace(r)) {
				roleMatch = true
				break
			}
		}

		if !roleMatch {
			// fmt.Printf("[DEBUG] RoleMiddleware: Access DENIED for role %s. Required: %v\n", roleStr, roles)
			utils.ErrorResponse(c, http.StatusForbidden, fmt.Sprintf("anda tidak memiliki akses (%s)", roleStr))
			c.Abort()
			return
		}

		// fmt.Printf("[DEBUG] RoleMiddleware: Access GRANTED for role %s\n", roleStr)
		c.Next()
	}
}
