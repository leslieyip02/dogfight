package room

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	secret []byte
}

func NewSession(secret []byte) *Session {
	return &Session{
		secret: secret,
	}
}

func (m *Session) createToken(roomId string, clientId string) (string, error) {
	claims := jwt.MapClaims{
		"roomId":   roomId,
		"clientId": clientId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Session) validateToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
