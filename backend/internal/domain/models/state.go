package models

import "time"

// GameState represents the current state of the game board
type GameState struct {
	ID                int64              `json:"id"`
	GameID            int64              `json:"game_id"`
	AvailableGems     map[string]int     `json:"available_gems"`
	VisibleCardsTier1 []DevelopmentCard  `json:"visible_cards_tier1"`
	VisibleCardsTier2 []DevelopmentCard  `json:"visible_cards_tier2"`
	VisibleCardsTier3 []DevelopmentCard  `json:"visible_cards_tier3"`
	AvailableNobles   []Noble            `json:"available_nobles"`
	DeckTier1         []DevelopmentCard  `json:"-"` // Hidden from clients
	DeckTier2         []DevelopmentCard  `json:"-"` // Hidden from clients
	DeckTier3         []DevelopmentCard  `json:"-"` // Hidden from clients
	DeckTier1Count    int                `json:"deck_tier1_count"`
	DeckTier2Count    int                `json:"deck_tier2_count"`
	DeckTier3Count    int                `json:"deck_tier3_count"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

// PlayerState represents a player's current resources and cards
type PlayerState struct {
	ID             int64              `json:"id"`
	GamePlayerID   int64              `json:"game_player_id"`
	Gems           map[string]int     `json:"gems"`
	PermanentGems  map[string]int     `json:"permanent_gems"`
	PurchasedCards []DevelopmentCard  `json:"purchased_cards"`
	ReservedCards  []DevelopmentCard  `json:"reserved_cards"`
	Nobles         []Noble            `json:"nobles"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// DevelopmentCard represents a development card
type DevelopmentCard struct {
	ID            int64          `json:"id"`
	Tier          int            `json:"tier"`
	GemType       string         `json:"gem_type"`
	VictoryPoints int            `json:"victory_points"`
	Cost          map[string]int `json:"cost"`
}

// Noble represents a noble
type Noble struct {
	ID            int64          `json:"id"`
	Name          string         `json:"name"`
	VictoryPoints int            `json:"victory_points"`
	Required      map[string]int `json:"required"`
}

// FullGameState contains all game information
type FullGameState struct {
	Game         *Game                     `json:"game"`
	Players      []*GamePlayer             `json:"players"`
	GameState    *GameState                `json:"game_state"`
	PlayerStates map[int64]*PlayerState    `json:"player_states"` // Keyed by user_id
}
