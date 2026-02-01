package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"splendor-backend/pkg/websocket"

	"github.com/gin-gonic/gin"
)

type GameplayHandler struct {
	engine GameplayEngine
	hub    *websocket.Hub
}

type GameplayEngine interface {
	TakeGems(ctx context.Context, gameID, userID int64, gems map[string]int) error
	PurchaseCard(ctx context.Context, gameID, userID int64, cardID int64, fromReserve bool) error
	ReserveCard(ctx context.Context, gameID, userID int64, cardID int64, tier int) error
}

func NewGameplayHandler(engine GameplayEngine, hub *websocket.Hub) *GameplayHandler {
	return &GameplayHandler{
		engine: engine,
		hub:    hub,
	}
}

type TakeGemsRequest struct {
	Gems map[string]int `json:"gems" binding:"required"`
}

type PurchaseCardRequest struct {
	CardID      int64 `json:"card_id" binding:"required"`
	FromReserve bool  `json:"from_reserve"`
}

type ReserveCardRequest struct {
	CardID int64 `json:"card_id"` // 0 for blind reserve
	Tier   int   `json:"tier"`    // Required for blind reserve
}

// TakeGems handles take gems action
func (h *GameplayHandler) TakeGems(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	var req TakeGemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.engine.TakeGems(c.Request.Context(), gameID, userID.(int64), req.Gems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast game update to all connected clients
	h.broadcastGameUpdate(gameIDStr, "game_update", gin.H{
		"action": "take_gems",
		"user_id": userID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Gems taken successfully"})
}

// PurchaseCard handles purchase card action
func (h *GameplayHandler) PurchaseCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	var req PurchaseCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.engine.PurchaseCard(c.Request.Context(), gameID, userID.(int64), req.CardID, req.FromReserve); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast game update to all connected clients
	h.broadcastGameUpdate(gameIDStr, "game_update", gin.H{
		"action": "purchase_card",
		"user_id": userID,
		"card_id": req.CardID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Card purchased successfully"})
}

// ReserveCard handles reserve card action
func (h *GameplayHandler) ReserveCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	var req ReserveCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.engine.ReserveCard(c.Request.Context(), gameID, userID.(int64), req.CardID, req.Tier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Broadcast game update to all connected clients
	h.broadcastGameUpdate(gameIDStr, "game_update", gin.H{
		"action": "reserve_card",
		"user_id": userID,
		"card_id": req.CardID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Card reserved successfully"})
}

// broadcastGameUpdate broadcasts a game update message to all connected clients
func (h *GameplayHandler) broadcastGameUpdate(gameID string, msgType string, payload gin.H) {
	message := map[string]interface{}{
		"type":    msgType,
		"payload": payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.hub.BroadcastToGame(gameID, messageBytes)
}
