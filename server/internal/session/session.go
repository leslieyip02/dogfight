package session

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// CreateToken issues a JWT.
func CreateToken(
	clientId string,
	username string,
	roomId string,
	secret []byte,
) (string, error) {
	claims := jwt.MapClaims{
		"clientId": clientId,
		"username": username,
		"roomId":   roomId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken returns the claims for a given tokenString.
func ParseToken(
	token string,
	secret []byte,
) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return t.Claims.(jwt.MapClaims), nil
}
