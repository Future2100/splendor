package models

import "time"

type GameStatus string

const (
	GameStatusWaiting    GameStatus = "waiting"
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusCompleted  GameStatus = "completed"
)

type Game struct {
	ID                 int64      `json:"id"`
	RoomCode           string     `json:"room_code"`
	Status             GameStatus `json:"status"`
	CurrentTurnPlayerID *int64    `json:"current_turn_player_id,omitempty"`
	TurnNumber         int        `json:"turn_number"`
	WinnerID           *int64     `json:"winner_id,omitempty"`
	CreatedBy          int64      `json:"created_by"`
	NumPlayers         int        `json:"num_players"`
	CreatedAt          time.Time  `json:"created_at"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`

	// Populated fields (not in DB)
	Players []*GamePlayer `json:"players,omitempty"`
	Creator *User         `json:"creator,omitempty"`
}

type GamePlayer struct {
	ID             int64     `json:"id"`
	GameID         int64     `json:"game_id"`
	UserID         int64     `json:"user_id"`
	PlayerPosition int       `json:"player_position"`
	VictoryPoints  int       `json:"victory_points"`
	IsActive       bool      `json:"is_active"`
	JoinedAt       time.Time `json:"joined_at"`

	// Populated field
	User *User `json:"user,omitempty"`
}

type CreateGameRequest struct {
	NumPlayers int `json:"num_players" binding:"required,min=2,max=4"`
}

type CreateGameResponse struct {
	Game     *Game  `json:"game"`
	RoomCode string `json:"room_code"`
}

type JoinGameRequest struct {
	RoomCode string `json:"room_code" binding:"required"`
}

type GameListResponse struct {
	Games []*Game `json:"games"`
	Total int     `json:"total"`
}
