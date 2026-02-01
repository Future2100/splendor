package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"splendor-backend/pkg/jwt"
	wshub "splendor-backend/pkg/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (configure based on your needs)
		return true
	},
}

type WebSocketHandler struct {
	hub       *wshub.Hub
	jwtSecret string
}

func NewWebSocketHandler(hub *wshub.Hub, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{
		hub:       hub,
		jwtSecret: jwtSecret,
	}
}

// HandleConnection handles WebSocket connection for a game
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	// Get game ID from URL
	gameIDStr := c.Param("id")
	_, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	// Get token from query param (for WebSocket auth)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Validate token
	claims, err := jwt.ValidateToken(token, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Create client
	client := &wshub.Client{
		ID:     strconv.FormatInt(claims.UserID, 10),
		GameID: gameIDStr,
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}

	// Register client
	h.hub.RegisterClient(client)

	// Start goroutines
	go h.writePump(client, conn)
	go h.readPump(client, conn)

	// Send welcome message
	welcomeMsg := wshub.Message{
		Type: "connected",
		Payload: map[string]interface{}{
			"message": "Connected to game",
			"game_id": gameIDStr,
			"user_id": claims.UserID,
		},
	}
	msgBytes, _ := json.Marshal(welcomeMsg)
	client.Send <- msgBytes
}

func (h *WebSocketHandler) readPump(client *wshub.Client, conn *websocket.Conn) {
	defer func() {
		h.hub.UnregisterClient(client)
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages
		var msg wshub.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Process message based on type
		h.handleMessage(client, &msg)
	}
}

func (h *WebSocketHandler) writePump(client *wshub.Client, conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) handleMessage(client *wshub.Client, msg *wshub.Message) {
	log.Printf("Received message from user %d: type=%s", client.UserID, msg.Type)

	// Handle different message types
	switch msg.Type {
	case "move":
		// Game move messages will be handled in Phase 6+
		log.Printf("Move message received: %+v", msg.Payload)

	case "chat":
		// Chat messages (if implemented)
		h.broadcastToGame(client.GameID, msg)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (h *WebSocketHandler) broadcastToGame(gameID string, msg *wshub.Message) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	h.hub.BroadcastToGame(gameID, msgBytes)
}
