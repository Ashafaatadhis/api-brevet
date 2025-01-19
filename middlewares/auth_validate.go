package middlewares

import (
	"new-brevet-be/config"
	"new-brevet-be/models"
	"new-brevet-be/services"
	"new-brevet-be/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// User struct type untuk token
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Level    int    `json:"level"`
	Token    string `json:"token"`
}

// AuthMiddleware untuk melindungi route dengan memverifikasi token JWT
func AuthMiddleware() fiber.Handler {
	db := config.DB
	return func(c *fiber.Ctx) error {
		var blacklist models.TokenBlacklist
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Jika token tidak ada, kirimkan response Unauthorized
			return utils.Response(c, fiber.StatusUnauthorized, "Token not provided", nil, nil, nil)
		}

		// Mengambil token dari format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Pastikan token diawali dengan "Bearer "
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid token format", nil, nil, nil)
		}

		// Memverifikasi token dan mendekodeklaimnya
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Kunci rahasia yang digunakan untuk memverifikasi tanda tangan JWT
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			// Jika token tidak valid atau error saat parsing, kirimkan Unauthorized
			return utils.Response(c, fiber.StatusUnauthorized, "Invalid or expired token", nil, nil, nil)
		}

		if err := services.CheckTokenBlacklist(db, tokenString, &blacklist); err != nil {
			return err
		}

		user := User{
			ID:       int(claims["sub"].(float64)), // Convert from float64 to uint
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
			Level:    int(claims["level"].(float64)),
			Token:    tokenString,
		}

		// Tambahkan informasi user ke dalam locals Fiber
		c.Locals("user", user)

		// Lanjutkan ke handler berikutnya
		return c.Next()
	}
}
