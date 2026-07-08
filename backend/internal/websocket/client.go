package websocket

import (
	"bytes"
	"log"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub       *Hub
	conn      *ws.Conn
	userID    string
	userEmail string
	send      chan []byte
}

func NewClient(hub *Hub, conn *ws.Conn, userID string, userEmail string) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		userID:    userID,
		userEmail: userEmail,
		send:      make(chan []byte, 256),
	}
}

func (client *Client) ReadPump() {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("unexpected websocket close error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		client.send <- message
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				client.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			writer, err := client.conn.NextWriter(ws.TextMessage)
			if err != nil {
				return
			}

			writer.Write(message)

			queuedMessages := len(client.send)
			for i := 0; i < queuedMessages; i++ {
				writer.Write([]byte("\n"))
				writer.Write(<-client.send)
			}

			if err := writer.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
