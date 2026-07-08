package websocket

import "fmt"

type Message struct {
	UserID string
	Data   []byte
}

type Hub struct {
	clients    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			if hub.clients[client.userID] == nil {
				hub.clients[client.userID] = make(map[*Client]bool)
			}

			hub.clients[client.userID][client] = true

			welcomeMessage := fmt.Sprintf(
				`{"type":"connected","user_id":"%s","message":"WebSocket connected"}`,
				client.userID,
			)

			client.send <- []byte(welcomeMessage)

		case client := <-hub.unregister:
			userClients, exists := hub.clients[client.userID]
			if !exists {
				continue
			}

			if _, ok := userClients[client]; ok {
				delete(userClients, client)
				close(client.send)
			}

			if len(userClients) == 0 {
				delete(hub.clients, client.userID)
			}

		case message := <-hub.broadcast:
			userClients, exists := hub.clients[message.UserID]
			if !exists {
				continue
			}

			for client := range userClients {
				select {
				case client.send <- message.Data:
				default:
					delete(userClients, client)
					close(client.send)
				}
			}
		}
	}
}

func (hub *Hub) Register(client *Client) {
	hub.register <- client
}

func (hub *Hub) SendToUser(userID string, data []byte) {
	hub.broadcast <- Message{
		UserID: userID,
		Data:   data,
	}
}
