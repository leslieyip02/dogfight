package room

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// A Session is responsible for managing tokens and client data.
type Session struct {
	secret   []byte
	upgrader websocket.Upgrader
}

type SessionClaims struct {
	roomId   string
	clientId string
	username string
}

func NewSession(secret []byte) *Session {
	return &Session{
		secret: secret,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// createToken issues a JWT.
func (s *Session) createToken(
	roomId string,
	clientId string,
	username string,
) (string, error) {
	claims := jwt.MapClaims{
		"roomId":   roomId,
		"clientId": clientId,
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// parseToken returns the claims for a given tokenString.
func (s *Session) parseToken(tokenString string) (*SessionClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse token %v", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return s.parseClaims(token.Claims.(jwt.MapClaims))
}

// parseClaims casts an untyped map into a typed map of SessionClaims.
func (s *Session) parseClaims(claims jwt.MapClaims) (*SessionClaims, error) {
	roomId, found := claims["roomId"].(string)
	if !found {
		return nil, fmt.Errorf("missing room ID")
	}

	clientId, found := claims["clientId"].(string)
	if !found {
		return nil, fmt.Errorf("missing client ID")
	}

	username, found := claims["username"].(string)
	if !found {
		username = "testificate"
	}

	return &SessionClaims{
		roomId:   roomId,
		clientId: clientId,
		username: username,
	}, nil
}

// createConn creates a WebSocket connection for the client.
func (s *Session) createConn(
	w *http.ResponseWriter,
	r *http.Request,
	clientId string,
	room *Room,
) (*websocket.Conn, error) {
	conn, err := s.upgrader.Upgrade(*w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection")
	}

	conn.SetCloseHandler(func(code int, text string) error {
		room.remove(clientId)
		return nil
	})

	return conn, nil
}
