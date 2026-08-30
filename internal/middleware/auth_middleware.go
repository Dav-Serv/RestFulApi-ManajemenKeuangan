package middleware

import (
	"RestFulApi-ManajemenKeuangan/config"
	"RestFulApi-ManajemenKeuangan/internal/utils"
	
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthRequired memvalidasi header "Authorization: Bearer <token>".
// Jika valid, user_id & email disimpan di gin.Context agar bisa dipakai handler.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.Error(c, http.StatusUnauthorized, "header Authorization tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "format header Authorization harus 'Bearer <token>'")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], config.App.JWTSecret)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "token tidak valid atau sudah kadaluwarsa")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}