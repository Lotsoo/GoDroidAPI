// file: websocket/hub.go
package websocket

import "log"

// Client adalah representasi dari koneksi websocket.
type Client struct {
	Hub  *Hub
	Conn *Conn // Wrapper untuk koneksi websocket
}

// Hub mengelola semua client dan pesan broadcast.
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Client registered, total clients:", len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Conn.Send)
				log.Println("Client unregistered, total clients:", len(h.Clients))
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Conn.Send <- message:
				default:
					close(client.Conn.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
