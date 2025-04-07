package entity

import (
	"abysslib/jwt"
	"github.com/fasthttp/websocket"
)

type Hub interface {
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
}

type Client struct {
	Hub            Hub
	authentication jwt.AuthenticationData
	conn           *websocket.Conn
	Send           chan []byte
}

func NewClient(
	hub Hub,
	authentication jwt.AuthenticationData,
	conn *websocket.Conn,
) *Client {
	return &Client{
		Hub:            hub,
		authentication: authentication,
		conn:           conn,
		Send:           make(chan []byte, 256),
	}
}

func (c *Client) GetAuthentication() jwt.AuthenticationData {
	return c.authentication
}

func (c *Client) CloseClient() error {
	close(c.Send)
	return c.conn.Close()
}

func (c *Client) WritePump() {
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(c.conn)

	for {
		message, ok := <-c.Send
		if !ok {
			// TODO: log
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			c.Hub.UnregisterClient(c)
			return
		}
	}
}
