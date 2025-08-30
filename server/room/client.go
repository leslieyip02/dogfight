package room

import (
	"server/game"
	"server/utils"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	username string
	conn     *websocket.Conn
	mu       sync.Mutex
}

func NewClient(username string) (*Client, error) {
	id, err := utils.NewShortId()
	if err != nil {
		return nil, err
	}

	client := Client{
		id:       id,
		username: username,
	}
	return &client, nil
}

func (c *Client) readPump(incoming chan<- []byte) {
	defer func() {
		quitEventMessage, err := game.NewQuitEventMessage(c.id)
		if err == nil {
			incoming <- quitEventMessage
		}
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		incoming <- message
	}
}

func (c *Client) writePump(outgoing <-chan []byte) {
	defer c.conn.Close()

	for data := range outgoing {
		err := c.writeMessage(data)
		if err != nil {
			break
		}
	}
}

func (c *Client) writeMessage(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, data)
}
