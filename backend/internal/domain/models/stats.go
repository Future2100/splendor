package models

import "time"

type GameStatistics struct {
	ID                   int64     `json:"id"`
	UserID               int64     `json:"user_id"`
	TotalGames           int       `json:"total_games"`
	TotalWins            int       `json:"total_wins"`
	TotalLosses          int       `json:"total_losses"`
	AveragePoints        float64   `json:"average_points"`
	AverageMovesPerGame  float64   `json:"average_moves_per_game"`
	FavoriteGemType      string    `json:"favorite_gem_type"`
	TotalNoblesEarned    int       `json:"total_nobles_earned"`
	TotalCardsPurchased  int       `json:"total_cards_purchased"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	UserID        int64   `json:"user_id"`
	Username      string  `json:"username"`
	TotalGames    int     `json:"total_games"`
	TotalWins     int     `json:"total_wins"`
	AveragePoints float64 `json:"average_points"`
	WinRate       float64 `json:"win_rate"`
}
