package postgres

import (
	"context"
	"fmt"

	"splendor-backend/internal/domain/models"
	"splendor-backend/pkg/database"
)

type CardRepository struct {
	db *database.DB
}

func NewCardRepository(db *database.DB) *CardRepository {
	return &CardRepository{db: db}
}

// GetAllCards retrieves all development cards
func (r *CardRepository) GetAllCards(ctx context.Context) ([]models.DevelopmentCard, error) {
	query := `
		SELECT id, tier, gem_type, victory_points,
		       cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx
		FROM development_cards
		ORDER BY tier, id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	cards := []models.DevelopmentCard{}
	for rows.Next() {
		card := models.DevelopmentCard{Cost: make(map[string]int)}
		var costDiamond, costSapphire, costEmerald, costRuby, costOnyx int

		err := rows.Scan(
			&card.ID,
			&card.Tier,
			&card.GemType,
			&card.VictoryPoints,
			&costDiamond,
			&costSapphire,
			&costEmerald,
			&costRuby,
			&costOnyx,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}

		card.Cost = map[string]int{
			"diamond":  costDiamond,
			"sapphire": costSapphire,
			"emerald":  costEmerald,
			"ruby":     costRuby,
			"onyx":     costOnyx,
		}

		cards = append(cards, card)
	}

	return cards, nil
}

// GetCardByID retrieves a card by ID
func (r *CardRepository) GetCardByID(ctx context.Context, id int64) (*models.DevelopmentCard, error) {
	query := `
		SELECT id, tier, gem_type, victory_points,
		       cost_diamond, cost_sapphire, cost_emerald, cost_ruby, cost_onyx
		FROM development_cards
		WHERE id = $1
	`

	card := &models.DevelopmentCard{Cost: make(map[string]int)}
	var costDiamond, costSapphire, costEmerald, costRuby, costOnyx int

	err := r.db.QueryRow(ctx, query, id).Scan(
		&card.ID,
		&card.Tier,
		&card.GemType,
		&card.VictoryPoints,
		&costDiamond,
		&costSapphire,
		&costEmerald,
		&costRuby,
		&costOnyx,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	card.Cost = map[string]int{
		"diamond":  costDiamond,
		"sapphire": costSapphire,
		"emerald":  costEmerald,
		"ruby":     costRuby,
		"onyx":     costOnyx,
	}

	return card, nil
}

// GetAllNobles retrieves all nobles
func (r *CardRepository) GetAllNobles(ctx context.Context) ([]models.Noble, error) {
	query := `
		SELECT id, name, victory_points,
		       required_diamond, required_sapphire, required_emerald,
		       required_ruby, required_onyx
		FROM nobles
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query nobles: %w", err)
	}
	defer rows.Close()

	nobles := []models.Noble{}
	for rows.Next() {
		noble := models.Noble{Required: make(map[string]int)}
		var reqDiamond, reqSapphire, reqEmerald, reqRuby, reqOnyx int

		err := rows.Scan(
			&noble.ID,
			&noble.Name,
			&noble.VictoryPoints,
			&reqDiamond,
			&reqSapphire,
			&reqEmerald,
			&reqRuby,
			&reqOnyx,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan noble: %w", err)
		}

		noble.Required = map[string]int{
			"diamond":  reqDiamond,
			"sapphire": reqSapphire,
			"emerald":  reqEmerald,
			"ruby":     reqRuby,
			"onyx":     reqOnyx,
		}

		nobles = append(nobles, noble)
	}

	return nobles, nil
}

// GetNobleByID retrieves a noble by ID
func (r *CardRepository) GetNobleByID(ctx context.Context, id int64) (*models.Noble, error) {
	query := `
		SELECT id, name, victory_points,
		       required_diamond, required_sapphire, required_emerald,
		       required_ruby, required_onyx
		FROM nobles
		WHERE id = $1
	`

	noble := &models.Noble{Required: make(map[string]int)}
	var reqDiamond, reqSapphire, reqEmerald, reqRuby, reqOnyx int

	err := r.db.QueryRow(ctx, query, id).Scan(
		&noble.ID,
		&noble.Name,
		&noble.VictoryPoints,
		&reqDiamond,
		&reqSapphire,
		&reqEmerald,
		&reqRuby,
		&reqOnyx,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get noble: %w", err)
	}

	noble.Required = map[string]int{
		"diamond":  reqDiamond,
		"sapphire": reqSapphire,
		"emerald":  reqEmerald,
		"ruby":     reqRuby,
		"onyx":     reqOnyx,
	}

	return noble, nil
}
