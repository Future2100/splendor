package handlers

import (
	"net/http"
	"strconv"

	"splendor-backend/internal/domain/models"
	"splendor-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame creates a new game
func (h *GameHandler) CreateGame(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.gameService.CreateGame(c.Request.Context(), userID.(int64), req.NumPlayers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListGames lists all games
func (h *GameHandler) ListGames(c *gin.Context) {
	// Get query params
	statusStr := c.Query("status")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var status *models.GameStatus
	if statusStr != "" {
		s := models.GameStatus(statusStr)
		status = &s
	}

	resp, err := h.gameService.ListGames(c.Request.Context(), status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list games"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGame retrieves a game by ID
func (h *GameHandler) GetGame(c *gin.Context) {
	gameIDStr := c.Param("id")
	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	game, err := h.gameService.GetGameByID(c.Request.Context(), gameID)
	if err != nil {
		if err == service.ErrGameNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get game"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game})
}

// JoinGame allows a user to join a game
func (h *GameHandler) JoinGame(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := h.gameService.JoinGame(c.Request.Context(), userID.(int64), req.RoomCode)
	if err != nil {
		switch err {
		case service.ErrGameNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		case service.ErrGameFull:
			c.JSON(http.StatusConflict, gin.H{"error": "Game is full"})
		case service.ErrGameStarted:
			c.JSON(http.StatusConflict, gin.H{"error": "Game has already started"})
		case service.ErrAlreadyInGame:
			c.JSON(http.StatusConflict, gin.H{"error": "You are already in this game"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join game"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game})
}

// LeaveGame removes a player from a game
func (h *GameHandler) LeaveGame(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameIDStr := c.Param("id")

	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	err = h.gameService.LeaveGame(c.Request.Context(), gameID, userID.(int64))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left game successfully"})
}

// StartGame starts a game
func (h *GameHandler) StartGame(c *gin.Context) {
	userID, _ := c.Get("userID")
	gameIDStr := c.Param("id")

	gameID, err := strconv.ParseInt(gameIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	game, err := h.gameService.StartGame(c.Request.Context(), gameID, userID.(int64))
	if err != nil {
		switch err {
		case service.ErrNotGameCreator:
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the game creator can start the game"})
		case service.ErrGameStarted:
			c.JSON(http.StatusConflict, gin.H{"error": "Game has already started"})
		case service.ErrNotEnoughPlayers:
			c.JSON(http.StatusConflict, gin.H{"error": "Need at least 2 players to start"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start game"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": game})
}
