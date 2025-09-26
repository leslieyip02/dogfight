package room

import (
	"server/game"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	username string
	conn     *websocket.Conn
	send     chan []byte
	mu       sync.Mutex
}

func NewClient(id string, username string) *Client {
	return &Client{
		id:       id,
		username: username,
		send:     make(chan []byte),
		mu:       sync.Mutex{},
	}
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

func (c *Client) writePump() {
	defer c.conn.Close()
	for data := range c.send {
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
