package room

import (
	"sync"

	"github.com/gorilla/websocket"
)

// A Client manages the interaction between the user and the server.
type Client struct {
	id       string
	username string
	conn     *websocket.Conn
	send     chan []byte
	mu       sync.Mutex
}

func newClient(id string, username string, conn *websocket.Conn) *Client {
	return &Client{
		id:       id,
		username: username,
		send:     make(chan []byte),
		mu:       sync.Mutex{},
		conn:     conn,
	}
}

// readPump relays messages from the client to the incoming channel.
func (c *Client) readPump(incoming chan<- []byte) {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		incoming <- message
	}
}

// writePump relays messages from the room to the client.
func (c *Client) writePump() {
	defer c.conn.Close()

	for data := range c.send {
		err := c.writeMessage(data)
		if err != nil {
			break
		}
	}
}

// writeMessage writes binary messages to the client's conn.
func (c *Client) writeMessage(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
