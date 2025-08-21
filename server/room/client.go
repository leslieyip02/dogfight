package room

import (
	"server/game"
	"server/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	username string
	conn     *websocket.Conn
}

func NewClient(username string) (*Client, error) {
	id, err := utils.GetShortId()
	if err != nil {
		return nil, err
	}

	client := Client{
		id:       id,
		username: username,
	}
	return &client, nil
}

func (c *Client) readPump(send chan<- []byte) {
	defer func() {
		quitEventMessage, err := game.NewQuitEventMessage(c.id)
		if err == nil {
			send <- quitEventMessage
		}
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		send <- message
	}
}

func (c *Client) writePump(broadcast <-chan []byte) {
	defer c.conn.Close()

	for data := range broadcast {
		err := c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			break
		}
	}
}
