package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"splendor-backend/internal/domain/models"
	"splendor-backend/pkg/database"
)

type StateRepository struct {
	db *database.DB
}

func NewStateRepository(db *database.DB) *StateRepository {
	return &StateRepository{db: db}
}

// CreateGameState creates a new game state
func (r *StateRepository) CreateGameState(ctx context.Context, state *models.GameState) error {
	// Convert cards and nobles to JSON
	tier1JSON, _ := json.Marshal(state.VisibleCardsTier1)
	tier2JSON, _ := json.Marshal(state.VisibleCardsTier2)
	tier3JSON, _ := json.Marshal(state.VisibleCardsTier3)
	noblesJSON, _ := json.Marshal(state.AvailableNobles)
	gemsJSON, _ := json.Marshal(state.AvailableGems)
	deck1JSON, _ := json.Marshal(state.DeckTier1)
	deck2JSON, _ := json.Marshal(state.DeckTier2)
	deck3JSON, _ := json.Marshal(state.DeckTier3)

	query := `
		INSERT INTO game_state (
			game_id, available_gems,
			visible_cards_tier1, visible_cards_tier2, visible_cards_tier3,
			available_nobles,
			deck_tier1, deck_tier2, deck_tier3,
			deck_tier1_count, deck_tier2_count, deck_tier3_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, updated_at
	`

	state.DeckTier1Count = len(state.DeckTier1)
	state.DeckTier2Count = len(state.DeckTier2)
	state.DeckTier3Count = len(state.DeckTier3)

	err := r.db.QueryRow(ctx, query,
		state.GameID,
		gemsJSON,
		tier1JSON,
		tier2JSON,
		tier3JSON,
		noblesJSON,
		deck1JSON,
		deck2JSON,
		deck3JSON,
		state.DeckTier1Count,
		state.DeckTier2Count,
		state.DeckTier3Count,
	).Scan(&state.ID, &state.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create game state: %w", err)
	}

	return nil
}

// GetGameState retrieves the game state
func (r *StateRepository) GetGameState(ctx context.Context, gameID int64) (*models.GameState, error) {
	query := `
		SELECT id, game_id, available_gems,
		       visible_cards_tier1, visible_cards_tier2, visible_cards_tier3,
		       available_nobles,
		       deck_tier1, deck_tier2, deck_tier3,
		       deck_tier1_count, deck_tier2_count, deck_tier3_count,
		       updated_at
		FROM game_state
		WHERE game_id = $1
	`

	state := &models.GameState{}
	var gemsJSON, tier1JSON, tier2JSON, tier3JSON, noblesJSON []byte
	var deck1JSON, deck2JSON, deck3JSON []byte

	err := r.db.QueryRow(ctx, query, gameID).Scan(
		&state.ID,
		&state.GameID,
		&gemsJSON,
		&tier1JSON,
		&tier2JSON,
		&tier3JSON,
		&noblesJSON,
		&deck1JSON,
		&deck2JSON,
		&deck3JSON,
		&state.DeckTier1Count,
		&state.DeckTier2Count,
		&state.DeckTier3Count,
		&state.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get game state: %w", err)
	}

	// Unmarshal JSON
	json.Unmarshal(gemsJSON, &state.AvailableGems)
	json.Unmarshal(tier1JSON, &state.VisibleCardsTier1)
	json.Unmarshal(tier2JSON, &state.VisibleCardsTier2)
	json.Unmarshal(tier3JSON, &state.VisibleCardsTier3)
	json.Unmarshal(noblesJSON, &state.AvailableNobles)
	json.Unmarshal(deck1JSON, &state.DeckTier1)
	json.Unmarshal(deck2JSON, &state.DeckTier2)
	json.Unmarshal(deck3JSON, &state.DeckTier3)

	return state, nil
}

// UpdateGameState updates the game state
func (r *StateRepository) UpdateGameState(ctx context.Context, state *models.GameState) error {
	tier1JSON, _ := json.Marshal(state.VisibleCardsTier1)
	tier2JSON, _ := json.Marshal(state.VisibleCardsTier2)
	tier3JSON, _ := json.Marshal(state.VisibleCardsTier3)
	noblesJSON, _ := json.Marshal(state.AvailableNobles)
	gemsJSON, _ := json.Marshal(state.AvailableGems)
	deck1JSON, _ := json.Marshal(state.DeckTier1)
	deck2JSON, _ := json.Marshal(state.DeckTier2)
	deck3JSON, _ := json.Marshal(state.DeckTier3)

	query := `
		UPDATE game_state
		SET available_gems = $1,
		    visible_cards_tier1 = $2,
		    visible_cards_tier2 = $3,
		    visible_cards_tier3 = $4,
		    available_nobles = $5,
		    deck_tier1 = $6,
		    deck_tier2 = $7,
		    deck_tier3 = $8,
		    deck_tier1_count = $9,
		    deck_tier2_count = $10,
		    deck_tier3_count = $11
		WHERE game_id = $12
	`

	_, err := r.db.Exec(ctx, query,
		gemsJSON,
		tier1JSON,
		tier2JSON,
		tier3JSON,
		noblesJSON,
		deck1JSON,
		deck2JSON,
		deck3JSON,
		state.DeckTier1Count,
		state.DeckTier2Count,
		state.DeckTier3Count,
		state.GameID,
	)

	if err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	return nil
}

// CreatePlayerState creates a new player state
func (r *StateRepository) CreatePlayerState(ctx context.Context, state *models.PlayerState) error {
	purchasedCardsJSON, _ := json.Marshal(state.PurchasedCards)
	reservedCardsJSON, _ := json.Marshal(state.ReservedCards)
	noblesJSON, _ := json.Marshal(state.Nobles)

	query := `
		INSERT INTO player_state (
			game_player_id,
			gems_diamond, gems_sapphire, gems_emerald, gems_ruby, gems_onyx, gems_gold,
			permanent_diamond, permanent_sapphire, permanent_emerald, permanent_ruby, permanent_onyx,
			purchased_cards, reserved_cards, nobles
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		state.GamePlayerID,
		state.Gems["diamond"],
		state.Gems["sapphire"],
		state.Gems["emerald"],
		state.Gems["ruby"],
		state.Gems["onyx"],
		state.Gems["gold"],
		state.PermanentGems["diamond"],
		state.PermanentGems["sapphire"],
		state.PermanentGems["emerald"],
		state.PermanentGems["ruby"],
		state.PermanentGems["onyx"],
		purchasedCardsJSON,
		reservedCardsJSON,
		noblesJSON,
	).Scan(&state.ID, &state.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create player state: %w", err)
	}

	return nil
}

// GetPlayerState retrieves a player's state
func (r *StateRepository) GetPlayerState(ctx context.Context, gamePlayerID int64) (*models.PlayerState, error) {
	query := `
		SELECT id, game_player_id,
		       gems_diamond, gems_sapphire, gems_emerald, gems_ruby, gems_onyx, gems_gold,
		       permanent_diamond, permanent_sapphire, permanent_emerald, permanent_ruby, permanent_onyx,
		       purchased_cards, reserved_cards, nobles,
		       updated_at
		FROM player_state
		WHERE game_player_id = $1
	`

	state := &models.PlayerState{
		Gems:          make(map[string]int),
		PermanentGems: make(map[string]int),
	}
	var purchasedCardsJSON, reservedCardsJSON, noblesJSON []byte
	var diamond, sapphire, emerald, ruby, onyx, gold int
	var permDiamond, permSapphire, permEmerald, permRuby, permOnyx int

	err := r.db.QueryRow(ctx, query, gamePlayerID).Scan(
		&state.ID,
		&state.GamePlayerID,
		&diamond,
		&sapphire,
		&emerald,
		&ruby,
		&onyx,
		&gold,
		&permDiamond,
		&permSapphire,
		&permEmerald,
		&permRuby,
		&permOnyx,
		&purchasedCardsJSON,
		&reservedCardsJSON,
		&noblesJSON,
		&state.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get player state: %w", err)
	}

	// Assign to maps - ensure all gem types exist even if 0
	state.Gems["diamond"] = diamond
	state.Gems["sapphire"] = sapphire
	state.Gems["emerald"] = emerald
	state.Gems["ruby"] = ruby
	state.Gems["onyx"] = onyx
	state.Gems["gold"] = gold

	state.PermanentGems["diamond"] = permDiamond
	state.PermanentGems["sapphire"] = permSapphire
	state.PermanentGems["emerald"] = permEmerald
	state.PermanentGems["ruby"] = permRuby
	state.PermanentGems["onyx"] = permOnyx

	// Unmarshal JSON
	json.Unmarshal(purchasedCardsJSON, &state.PurchasedCards)
	json.Unmarshal(reservedCardsJSON, &state.ReservedCards)
	json.Unmarshal(noblesJSON, &state.Nobles)

	return state, nil
}

// UpdatePlayerState updates a player's state
func (r *StateRepository) UpdatePlayerState(ctx context.Context, state *models.PlayerState) error {
	purchasedCardsJSON, _ := json.Marshal(state.PurchasedCards)
	reservedCardsJSON, _ := json.Marshal(state.ReservedCards)
	noblesJSON, _ := json.Marshal(state.Nobles)

	query := `
		UPDATE player_state
		SET gems_diamond = $1, gems_sapphire = $2, gems_emerald = $3,
		    gems_ruby = $4, gems_onyx = $5, gems_gold = $6,
		    permanent_diamond = $7, permanent_sapphire = $8, permanent_emerald = $9,
		    permanent_ruby = $10, permanent_onyx = $11,
		    purchased_cards = $12, reserved_cards = $13, nobles = $14
		WHERE game_player_id = $15
	`

	_, err := r.db.Exec(ctx, query,
		state.Gems["diamond"],
		state.Gems["sapphire"],
		state.Gems["emerald"],
		state.Gems["ruby"],
		state.Gems["onyx"],
		state.Gems["gold"],
		state.PermanentGems["diamond"],
		state.PermanentGems["sapphire"],
		state.PermanentGems["emerald"],
		state.PermanentGems["ruby"],
		state.PermanentGems["onyx"],
		purchasedCardsJSON,
		reservedCardsJSON,
		noblesJSON,
		state.GamePlayerID,
	)

	if err != nil {
		return fmt.Errorf("failed to update player state: %w", err)
	}

	return nil
}
