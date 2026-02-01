package postgres

import (
	"context"
	"fmt"

	"splendor-backend/internal/domain/models"
	"splendor-backend/pkg/database"
)

type StatsRepository struct {
	db *database.DB
}

func NewStatsRepository(db *database.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetUserStats retrieves user statistics
func (r *StatsRepository) GetUserStats(ctx context.Context, userID int64) (*models.GameStatistics, error) {
	query := `
		SELECT id, user_id, total_games, total_wins, total_losses,
		       average_points, average_moves_per_game, favorite_gem_type,
		       total_nobles_earned, total_cards_purchased, updated_at
		FROM game_statistics
		WHERE user_id = $1
	`

	stats := &models.GameStatistics{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stats.ID,
		&stats.UserID,
		&stats.TotalGames,
		&stats.TotalWins,
		&stats.TotalLosses,
		&stats.AveragePoints,
		&stats.AverageMovesPerGame,
		&stats.FavoriteGemType,
		&stats.TotalNoblesEarned,
		&stats.TotalCardsPurchased,
		&stats.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("stats not found")
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}

// GetLeaderboard retrieves the leaderboard
func (r *StatsRepository) GetLeaderboard(ctx context.Context, limit, offset int) ([]*models.LeaderboardEntry, error) {
	query := `
		SELECT gs.user_id, u.username, gs.total_games, gs.total_wins,
		       gs.average_points, CAST(gs.total_wins AS FLOAT) / NULLIF(gs.total_games, 0) as win_rate
		FROM game_statistics gs
		JOIN users u ON u.id = gs.user_id
		WHERE gs.total_games > 0
		ORDER BY win_rate DESC, gs.total_wins DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer rows.Close()

	entries := []*models.LeaderboardEntry{}
	for rows.Next() {
		entry := &models.LeaderboardEntry{}
		err := rows.Scan(
			&entry.UserID,
			&entry.Username,
			&entry.TotalGames,
			&entry.TotalWins,
			&entry.AveragePoints,
			&entry.WinRate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// UpdateStats updates or creates user statistics
func (r *StatsRepository) UpdateStats(ctx context.Context, stats *models.GameStatistics) error {
	query := `
		INSERT INTO game_statistics (
			user_id, total_games, total_wins, total_losses,
			average_points, average_moves_per_game, favorite_gem_type,
			total_nobles_earned, total_cards_purchased
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id) DO UPDATE SET
			total_games = $2,
			total_wins = $3,
			total_losses = $4,
			average_points = $5,
			average_moves_per_game = $6,
			favorite_gem_type = $7,
			total_nobles_earned = $8,
			total_cards_purchased = $9
		RETURNING id, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		stats.UserID,
		stats.TotalGames,
		stats.TotalWins,
		stats.TotalLosses,
		stats.AveragePoints,
		stats.AverageMovesPerGame,
		stats.FavoriteGemType,
		stats.TotalNoblesEarned,
		stats.TotalCardsPurchased,
	).Scan(&stats.ID, &stats.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update stats: %w", err)
	}

	return nil
}
