package hub

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	"github.com/intezya/pkglib/logger"
)

const sendBufferSize = 256

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
		Send:           make(chan []byte, sendBufferSize),
		connectTime:    time.Now(),
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
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))
	c.conn.SetPongHandler(
		func(string) error {
			_ = c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))

			return nil
		},
	)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				logger.Log.Debugf("error: %v", err)
			}

			break
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(connectionTimeout))

		logger.Log.Debugf("Received message from user %d: %s", c.authentication.ID(), message)

		c.Send <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(connectionPingPeriod)
	defer func() {
		ticker.Stop()

		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWaitTimeout))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})

				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(message)

			n := len(c.Send)
			for range n {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWaitTimeout))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Log.Debugf("Error sending ping to user %d: %v", c.authentication.ID(), err)

				return
			}

			logger.Log.Debugf("Sent ping to user %d", c.authentication.ID())
		}
	}
}
