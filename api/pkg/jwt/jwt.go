package jwt

import (
	"time"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

const defaultTTL = 24 * time.Hour

// GenerateToken creates a signed JWT containing the user's id and username.
func GenerateToken(secret string, userID int64, username string) (string, error) {
	claims := golangjwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(defaultTTL).Unix(),
	}
	token := golangjwt.NewWithClaims(golangjwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
