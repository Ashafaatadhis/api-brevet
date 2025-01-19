package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SecretKey adalah key untuk signing JWT

// GenerateToken untuk membuat JWT token
func GenerateToken(userID int, username string, role string, level int) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	var timeExpiryToken = os.Getenv("TOKEN_EXPIRY")
	expiryInHours, err := strconv.Atoi(timeExpiryToken)
	if err != nil {
		fmt.Println("Error parsing TOKEN_EXPIRY:", err)
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["sub"] = userID
	claims["username"] = username
	claims["role"] = role
	claims["level"] = level
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expiryInHours)).Unix() // Expiry 24 jam

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken untuk memverifikasi dan membaca token JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return secretKey, nil
	})
	return token, err
}
