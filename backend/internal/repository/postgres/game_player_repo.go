package postgres

import (
	"context"
	"fmt"

	"splendor-backend/internal/domain/models"
)

// UpdatePlayer updates a game player's victory points
func (r *GameRepository) UpdatePlayer(ctx context.Context, player *models.GamePlayer) error {
	query := `
		UPDATE game_players
		SET victory_points = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, player.VictoryPoints, player.ID)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	return nil
}
