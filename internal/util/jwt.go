package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken builds a signed JWT with user and org context embedded.
func NewToken(userID, orgID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"org_id":  orgID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
