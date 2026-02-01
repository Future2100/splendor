package service

import (
	"context"

	"splendor-backend/internal/domain/models"
	"splendor-backend/internal/repository/postgres"
)

type StatsService struct {
	statsRepo *postgres.StatsRepository
}

func NewStatsService(statsRepo *postgres.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

// GetUserStats retrieves user statistics
func (s *StatsService) GetUserStats(ctx context.Context, userID int64) (*models.GameStatistics, error) {
	return s.statsRepo.GetUserStats(ctx, userID)
}

// GetLeaderboard retrieves the leaderboard
func (s *StatsService) GetLeaderboard(ctx context.Context, limit, offset int) ([]*models.LeaderboardEntry, error) {
	return s.statsRepo.GetLeaderboard(ctx, limit, offset)
}

// UpdateGameStats updates statistics after a game ends
func (s *StatsService) UpdateGameStats(ctx context.Context, game *models.Game, players []*models.GamePlayer) error {
	for _, player := range players {
		stats, err := s.statsRepo.GetUserStats(ctx, player.UserID)
		if err != nil {
			// Create new stats if not exists
			stats = &models.GameStatistics{
				UserID: player.UserID,
			}
		}

		stats.TotalGames++
		if game.WinnerID != nil && *game.WinnerID == player.UserID {
			stats.TotalWins++
		} else {
			stats.TotalLosses++
		}

		// Update average points
		stats.AveragePoints = float64(stats.TotalWins*player.VictoryPoints+int(stats.AveragePoints)*stats.TotalLosses) / float64(stats.TotalGames)

		if err := s.statsRepo.UpdateStats(ctx, stats); err != nil {
			return err
		}
	}

	return nil
}
