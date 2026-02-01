package websocket

import (
	"sync"
)

// Message represents a WebSocket message
type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Client represents a WebSocket client connection
type Client struct {
	ID     string
	GameID string
	UserID int64
	Conn   any // Will be *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients by game ID
	games map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to all clients in a game
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// BroadcastMessage contains a message and target game ID
type BroadcastMessage struct {
	GameID  string
	Message []byte
}

func NewHub() *Hub {
	return &Hub{
		games:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToGame(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.games[client.GameID]; !ok {
		h.games[client.GameID] = make(map[*Client]bool)
	}
	h.games[client.GameID][client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.games[client.GameID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)

			// Clean up empty game rooms
			if len(clients) == 0 {
				delete(h.games, client.GameID)
			}
		}
	}
}

func (h *Hub) broadcastToGame(message *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.games[message.GameID]; ok {
		for client := range clients {
			select {
			case client.Send <- message.Message:
			default:
				// Client send buffer is full, close it
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// RegisterClient registers a new client
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastToGame sends a message to all clients in a specific game
func (h *Hub) BroadcastToGame(gameID string, message []byte) {
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: message,
	}
}

// GetGameClientCount returns the number of connected clients for a game
func (h *Hub) GetGameClientCount(gameID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.games[gameID]; ok {
		return len(clients)
	}
	return 0
}
