package middlewares

import (
	"net/http"
	"strings"
	"tokobiru/config"
	"tokobiru/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah middleware untuk memeriksa token JWT.
// Middleware ini harus dijalankan pertama untuk rute yang diproteksi.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := config.LoadConfig()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not load configuration"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Memeriksa format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format, 'Bearer' prefix not found"})
			return
		}

		// Memanggil fungsi dari paket 'services' untuk validasi
		claims, err := services.ValidateToken(tokenString, cfg.JWTSecretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "details": err.Error()})
			return
		}

		// Menyimpan info user ke dalam context untuk digunakan oleh handler lain
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		// Lanjut ke handler/middleware berikutnya
		c.Next()
	}
}

// RoleMiddleware adalah middleware untuk memeriksa role user.
// Middleware ini harus dijalankan SETELAH AuthMiddleware.
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User role not found in context. Ensure AuthMiddleware runs first."})
			return
		}

		if userRole.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			return
		}

		c.Next()
	}
}
