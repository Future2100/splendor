package handlers

import (
	"context"
	"net/http"
	"strconv"

	"splendor-backend/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService StatsService
}

type StatsService interface {
	GetUserStats(ctx context.Context, userID int64) (*models.GameStatistics, error)
	GetLeaderboard(ctx context.Context, limit, offset int) ([]*models.LeaderboardEntry, error)
}

func NewStatsHandler(statsService StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetUserStats retrieves user statistics
func (h *StatsHandler) GetUserStats(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	stats, err := h.statsService.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stats not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetLeaderboard retrieves the leaderboard
func (h *StatsHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	leaderboard, err := h.statsService.GetLeaderboard(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}
