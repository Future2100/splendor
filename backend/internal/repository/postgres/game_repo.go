package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"splendor-backend/internal/domain/models"
	"splendor-backend/pkg/database"

	"github.com/jackc/pgx/v5"
)

type GameRepository struct {
	db *database.DB
}

func NewGameRepository(db *database.DB) *GameRepository {
	return &GameRepository{db: db}
}

// GenerateRoomCode generates a unique 6-character room code
func (r *GameRepository) GenerateRoomCode(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ { // Try 10 times
		code := generateRandomCode(6)

		// Check if code exists
		var exists bool
		err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM games WHERE room_code = $1)", code).Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique room code")
}

func generateRandomCode(length int) string {
	bytes := make([]byte, length/2+1)
	rand.Read(bytes)
	code := hex.EncodeToString(bytes)[:length]
	return code
}

// Create creates a new game
func (r *GameRepository) Create(ctx context.Context, game *models.Game) error {
	query := `
		INSERT INTO games (room_code, status, num_players, created_by, turn_number)
		VALUES ($1, $2, $3, $4, 0)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, game.RoomCode, game.Status, game.NumPlayers, game.CreatedBy).
		Scan(&game.ID, &game.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create game: %w", err)
	}

	return nil
}

// GetByID retrieves a game by ID
func (r *GameRepository) GetByID(ctx context.Context, id int64) (*models.Game, error) {
	query := `
		SELECT id, room_code, status, current_turn_player_id, turn_number,
		       winner_id, created_by, num_players, created_at, started_at, completed_at
		FROM games
		WHERE id = $1
	`

	game := &models.Game{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&game.ID,
		&game.RoomCode,
		&game.Status,
		&game.CurrentTurnPlayerID,
		&game.TurnNumber,
		&game.WinnerID,
		&game.CreatedBy,
		&game.NumPlayers,
		&game.CreatedAt,
		&game.StartedAt,
		&game.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("game not found")
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

// GetByRoomCode retrieves a game by room code
func (r *GameRepository) GetByRoomCode(ctx context.Context, roomCode string) (*models.Game, error) {
	query := `
		SELECT id, room_code, status, current_turn_player_id, turn_number,
		       winner_id, created_by, num_players, created_at, started_at, completed_at
		FROM games
		WHERE room_code = $1
	`

	game := &models.Game{}
	err := r.db.QueryRow(ctx, query, roomCode).Scan(
		&game.ID,
		&game.RoomCode,
		&game.Status,
		&game.CurrentTurnPlayerID,
		&game.TurnNumber,
		&game.WinnerID,
		&game.CreatedBy,
		&game.NumPlayers,
		&game.CreatedAt,
		&game.StartedAt,
		&game.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("game not found")
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

// List retrieves games with optional status filter
func (r *GameRepository) List(ctx context.Context, status *models.GameStatus, limit, offset int) ([]*models.Game, int, error) {
	// Count total
	var countQuery string
	var args []interface{}

	if status != nil {
		countQuery = "SELECT COUNT(*) FROM games WHERE status = $1"
		args = append(args, *status)
	} else {
		countQuery = "SELECT COUNT(*) FROM games"
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count games: %w", err)
	}

	// Get games
	var query string
	args = []interface{}{}

	if status != nil {
		query = `
			SELECT id, room_code, status, current_turn_player_id, turn_number,
			       winner_id, created_by, num_players, created_at, started_at, completed_at
			FROM games
			WHERE status = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`
		args = append(args, *status, limit, offset)
	} else {
		query = `
			SELECT id, room_code, status, current_turn_player_id, turn_number,
			       winner_id, created_by, num_players, created_at, started_at, completed_at
			FROM games
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list games: %w", err)
	}
	defer rows.Close()

	games := []*models.Game{}
	for rows.Next() {
		game := &models.Game{}
		err := rows.Scan(
			&game.ID,
			&game.RoomCode,
			&game.Status,
			&game.CurrentTurnPlayerID,
			&game.TurnNumber,
			&game.WinnerID,
			&game.CreatedBy,
			&game.NumPlayers,
			&game.CreatedAt,
			&game.StartedAt,
			&game.CompletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	return games, total, nil
}

// Update updates a game
func (r *GameRepository) Update(ctx context.Context, game *models.Game) error {
	query := `
		UPDATE games
		SET status = $1, current_turn_player_id = $2, turn_number = $3,
		    winner_id = $4, started_at = $5, completed_at = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(ctx, query,
		game.Status,
		game.CurrentTurnPlayerID,
		game.TurnNumber,
		game.WinnerID,
		game.StartedAt,
		game.CompletedAt,
		game.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}

// AddPlayer adds a player to a game
func (r *GameRepository) AddPlayer(ctx context.Context, gamePlayer *models.GamePlayer) error {
	query := `
		INSERT INTO game_players (game_id, user_id, player_position, victory_points, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, joined_at
	`

	err := r.db.QueryRow(ctx, query,
		gamePlayer.GameID,
		gamePlayer.UserID,
		gamePlayer.PlayerPosition,
		gamePlayer.VictoryPoints,
		gamePlayer.IsActive,
	).Scan(&gamePlayer.ID, &gamePlayer.JoinedAt)

	if err != nil {
		return fmt.Errorf("failed to add player: %w", err)
	}

	return nil
}

// GetPlayers retrieves all players for a game
func (r *GameRepository) GetPlayers(ctx context.Context, gameID int64) ([]*models.GamePlayer, error) {
	query := `
		SELECT gp.id, gp.game_id, gp.user_id, gp.player_position,
		       gp.victory_points, gp.is_active, gp.joined_at,
		       u.id, u.username, u.email, u.created_at, u.updated_at
		FROM game_players gp
		JOIN users u ON u.id = gp.user_id
		WHERE gp.game_id = $1
		ORDER BY gp.player_position
	`

	rows, err := r.db.Query(ctx, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}
	defer rows.Close()

	players := []*models.GamePlayer{}
	for rows.Next() {
		player := &models.GamePlayer{User: &models.User{}}
		err := rows.Scan(
			&player.ID,
			&player.GameID,
			&player.UserID,
			&player.PlayerPosition,
			&player.VictoryPoints,
			&player.IsActive,
			&player.JoinedAt,
			&player.User.ID,
			&player.User.Username,
			&player.User.Email,
			&player.User.CreatedAt,
			&player.User.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}
		players = append(players, player)
	}

	return players, nil
}

// GetPlayerCount returns the number of players in a game
func (r *GameRepository) GetPlayerCount(ctx context.Context, gameID int64) (int, error) {
	query := `SELECT COUNT(*) FROM game_players WHERE game_id = $1 AND is_active = true`

	var count int
	err := r.db.QueryRow(ctx, query, gameID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count players: %w", err)
	}

	return count, nil
}

// IsPlayerInGame checks if a user is already in a game
func (r *GameRepository) IsPlayerInGame(ctx context.Context, gameID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM game_players WHERE game_id = $1 AND user_id = $2 AND is_active = true)`

	var exists bool
	err := r.db.QueryRow(ctx, query, gameID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check player in game: %w", err)
	}

	return exists, nil
}

// RemovePlayer marks a player as inactive
func (r *GameRepository) RemovePlayer(ctx context.Context, gameID, userID int64) error {
	query := `UPDATE game_players SET is_active = false WHERE game_id = $1 AND user_id = $2`

	_, err := r.db.Exec(ctx, query, gameID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove player: %w", err)
	}

	return nil
}
