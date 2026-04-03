package ws

import (
	"sync"	 
	 "github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients map[uint]*websocket.Conn
	mu      sync.Mutex
}

var GlobalHub = &Hub{
	clients: make(map[uint]*websocket.Conn),
}

func (h *Hub) Register(userID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
}

func (h *Hub) Unregister(userID uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, userID)
}

func (h *Hub) Send(receiverID uint, message []byte) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	conn, online := h.clients[receiverID]
	if !online {
		return false
	}
	err := conn.WriteMessage(1, message)
	return err == nil
}

func (h *Hub) IsOnline(userID uint) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, online := h.clients[userID]
	return online 
}