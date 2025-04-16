package hub

import (
	"github.com/gorilla/websocket"
	"github.com/intezya/pkglib/logger"
	"time"
	"websocket/internal/domain/entity"
)

type Client struct {
	Hub            *Hub
	authentication *entity.AuthenticationData
	conn           *websocket.Conn
	Send           chan []byte
	connectTime    time.Time
}

func NewClient(
	hub *Hub,
	authentication *entity.AuthenticationData,
	conn *websocket.Conn,
) *Client {
	return &Client{
		Hub:            hub,
		authentication: authentication,
		conn:           conn,
		Send:           make(chan []byte, 256),
	}
}

func (c *Client) GetAuthentication() *entity.AuthenticationData {
	return c.authentication
}

func (c *Client) CloseClient() error {
	close(c.Send)
	return c.conn.Close()
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))
			return nil
		},
	)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Debugf("error: %v", err)
			}
			break
		}

		c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))

		logger.Log.Debugf("Received message from userentity %d: %s", c.authentication.ID(), message)

		c.Send <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(connectionPingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWaitTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWaitTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Log.Debugf("Error sending ping to userentity %d: %v", c.authentication.ID(), err)
				return
			}
			logger.Log.Debugf("Sent ping to userentity %d", c.authentication.ID())
		}
	}
}
