package handlers

import (
	"context"
	"net/http"
	"strconv"

	"splendor-backend/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type StateHandler struct {
	engine GameStateEngine
}

type GameStateEngine interface {
	GetGameState(ctx context.Context, gameID int64) (*models.FullGameState, error)
}

func NewStateHandler(engine GameStateEngine) *StateHandler {
	return &StateHandler{
		engine: engine,
	}
}

// GetGameState retrieves the full game state
func (h *StateHandler) GetGameState(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	state, err := h.engine.GetGameState(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": state})
}
