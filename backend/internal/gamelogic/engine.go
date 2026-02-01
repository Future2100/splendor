package gamelogic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"splendor-backend/internal/domain/models"
	"splendor-backend/internal/repository/postgres"
)

type GameEngine struct {
	gameRepo *postgres.GameRepository
	cardRepo *postgres.CardRepository
	stateRepo *postgres.StateRepository
}

func NewGameEngine(gameRepo *postgres.GameRepository, cardRepo *postgres.CardRepository, stateRepo *postgres.StateRepository) *GameEngine {
	return &GameEngine{
		gameRepo: gameRepo,
		cardRepo: cardRepo,
		stateRepo: stateRepo,
	}
}

// InitializeGame initializes a new game with shuffled cards, gems, and nobles
func (e *GameEngine) InitializeGame(ctx context.Context, gameID int64) error {
	_, err := e.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}

	players, err := e.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get players: %w", err)
	}

	numPlayers := len(players)

	// Get all development cards
	allCards, err := e.cardRepo.GetAllCards(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cards: %w", err)
	}

	// Separate cards by tier
	tier1Cards := []models.DevelopmentCard{}
	tier2Cards := []models.DevelopmentCard{}
	tier3Cards := []models.DevelopmentCard{}

	for _, card := range allCards {
		switch card.Tier {
		case 1:
			tier1Cards = append(tier1Cards, card)
		case 2:
			tier2Cards = append(tier2Cards, card)
		case 3:
			tier3Cards = append(tier3Cards, card)
		}
	}

	// Shuffle each tier
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffleCardsWithRNG(tier1Cards, rng)
	shuffleCardsWithRNG(tier2Cards, rng)
	shuffleCardsWithRNG(tier3Cards, rng)

	// Deal 4 cards from each tier
	visibleTier1 := tier1Cards[:4]
	visibleTier2 := tier2Cards[:4]
	visibleTier3 := tier3Cards[:4]

	deckTier1 := tier1Cards[4:]
	deckTier2 := tier2Cards[4:]
	deckTier3 := tier3Cards[4:]

	// Get all nobles and select numPlayers + 1
	allNobles, err := e.cardRepo.GetAllNobles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nobles: %w", err)
	}

	shuffleNoblesWithRNG(allNobles, rng)
	noblesCount := numPlayers + 1
	if noblesCount > len(allNobles) {
		noblesCount = len(allNobles)
	}
	selectedNobles := allNobles[:noblesCount]

	// Determine gem counts based on player count
	gemCounts := getGemCounts(numPlayers)

	// Create game state
	gameState := &models.GameState{
		GameID:        gameID,
		AvailableGems: gemCounts,
		VisibleCardsTier1: visibleTier1,
		VisibleCardsTier2: visibleTier2,
		VisibleCardsTier3: visibleTier3,
		AvailableNobles:  selectedNobles,
		DeckTier1:        deckTier1,
		DeckTier2:        deckTier2,
		DeckTier3:        deckTier3,
	}

	if err := e.stateRepo.CreateGameState(ctx, gameState); err != nil {
		return fmt.Errorf("failed to create game state: %w", err)
	}

	// Initialize player states
	for _, player := range players {
		playerState := &models.PlayerState{
			GamePlayerID: player.ID,
			Gems:         make(map[string]int),
			PermanentGems: make(map[string]int),
			PurchasedCards: []models.DevelopmentCard{},
			ReservedCards:  []models.DevelopmentCard{},
			Nobles:         []models.Noble{},
		}

		// Initialize all gem types to 0
		gemTypes := []string{"diamond", "sapphire", "emerald", "ruby", "onyx", "gold"}
		for _, gemType := range gemTypes {
			playerState.Gems[gemType] = 0
			playerState.PermanentGems[gemType] = 0
		}

		if err := e.stateRepo.CreatePlayerState(ctx, playerState); err != nil {
			return fmt.Errorf("failed to create player state: %w", err)
		}
	}

	return nil
}

func shuffleCardsWithRNG(cards []models.DevelopmentCard, rng *rand.Rand) {
	for i := range cards {
		j := rng.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func shuffleNoblesWithRNG(nobles []models.Noble, rng *rand.Rand) {
	for i := range nobles {
		j := rng.Intn(i + 1)
		nobles[i], nobles[j] = nobles[j], nobles[i]
	}
}

// getGemCounts returns gem counts based on player count
// Splendor rules: 2 players = 4 gems, 3 players = 5 gems, 4 players = 7 gems
// Gold coins: always 5
func getGemCounts(numPlayers int) map[string]int {
	baseCount := 0
	switch numPlayers {
	case 2:
		baseCount = 4
	case 3:
		baseCount = 5
	case 4:
		baseCount = 7
	default:
		baseCount = 7
	}

	return map[string]int{
		"diamond":  baseCount,
		"sapphire": baseCount,
		"emerald":  baseCount,
		"ruby":     baseCount,
		"onyx":     baseCount,
		"gold":     5, // Always 5 gold coins
	}
}

// GetGameState retrieves the full game state including player states
func (e *GameEngine) GetGameState(ctx context.Context, gameID int64) (*models.FullGameState, error) {
	game, err := e.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	players, err := e.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	gameState, err := e.stateRepo.GetGameState(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game state: %w", err)
	}

	playerStates := make(map[int64]*models.PlayerState)
	for _, player := range players {
		state, err := e.stateRepo.GetPlayerState(ctx, player.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get player state: %w", err)
		}
		playerStates[player.UserID] = state
	}

	return &models.FullGameState{
		Game:         game,
		Players:      players,
		GameState:    gameState,
		PlayerStates: playerStates,
	}, nil
}
