package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT issues a signed JWT for a user id, email and role.
func GenerateJWT(secret string, userID uint, email, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(ttl).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
