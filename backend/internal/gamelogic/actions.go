package gamelogic

import (
	"context"
	"fmt"
	"time"

	"splendor-backend/internal/domain/models"
)

// TakeGems executes the take gems action
func (e *GameEngine) TakeGems(ctx context.Context, gameID, userID int64, gems map[string]int) error {
	validator := NewGameValidator()

	// Get game
	game, err := e.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}

	// Validate turn
	if err := validator.ValidateTurn(game, userID); err != nil {
		return err
	}

	// Get game state
	gameState, err := e.stateRepo.GetGameState(ctx, gameID)
	if err != nil {
		return err
	}

	// Get player
	players, err := e.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return err
	}

	var currentPlayer *models.GamePlayer
	for _, p := range players {
		if p.UserID == userID {
			currentPlayer = p
			break
		}
	}

	// Get player state
	playerState, err := e.stateRepo.GetPlayerState(ctx, currentPlayer.ID)
	if err != nil {
		return err
	}

	// Validate action
	if err := validator.ValidateTakeGems(gameState, playerState, gems); err != nil {
		return err
	}

	// Execute: Remove gems from bank, add to player
	for gemType, count := range gems {
		if count > 0 {
			gameState.AvailableGems[gemType] -= count
			playerState.Gems[gemType] += count
		}
	}

	// Update game state
	if err := e.stateRepo.UpdateGameState(ctx, gameState); err != nil {
		return err
	}

	// Update player state
	if err := e.stateRepo.UpdatePlayerState(ctx, playerState); err != nil {
		return err
	}

	// Switch turn
	if err := e.switchTurn(ctx, game, players); err != nil {
		return err
	}

	return nil
}

// PurchaseCard executes the purchase card action
func (e *GameEngine) PurchaseCard(ctx context.Context, gameID, userID int64, cardID int64, fromReserve bool) error {
	validator := NewGameValidator()

	// Get game
	game, err := e.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}

	// Validate turn
	if err := validator.ValidateTurn(game, userID); err != nil {
		return err
	}

	// Get card
	card, err := e.cardRepo.GetCardByID(ctx, cardID)
	if err != nil {
		return err
	}

	// Get player
	players, err := e.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return err
	}

	var currentPlayer *models.GamePlayer
	for _, p := range players {
		if p.UserID == userID {
			currentPlayer = p
			break
		}
	}

	// Get states
	gameState, err := e.stateRepo.GetGameState(ctx, gameID)
	if err != nil {
		return err
	}

	playerState, err := e.stateRepo.GetPlayerState(ctx, currentPlayer.ID)
	if err != nil {
		return err
	}

	// Validate can afford
	if err := validator.ValidatePurchaseCard(gameState, playerState, card); err != nil {
		return err
	}

	// Calculate actual cost
	actualCost := validator.CalculateCost(card, playerState)

	// Pay cost
	for gemType, cost := range actualCost {
		playerState.Gems[gemType] -= cost
		gameState.AvailableGems[gemType] += cost
	}

	// Add card to player
	playerState.PurchasedCards = append(playerState.PurchasedCards, *card)
	playerState.PermanentGems[card.GemType]++

	// Remove from reserve if applicable
	if fromReserve {
		newReserved := []models.DevelopmentCard{}
		for _, c := range playerState.ReservedCards {
			if c.ID != cardID {
				newReserved = append(newReserved, c)
			}
		}
		playerState.ReservedCards = newReserved
	} else {
		// Remove from visible cards and replace
		e.removeAndReplaceCard(gameState, card)
	}

	// Update victory points
	currentPlayer.VictoryPoints += card.VictoryPoints
	if err := e.gameRepo.UpdatePlayer(ctx, currentPlayer); err != nil {
		return err
	}

	// Check for noble visits
	for i, noble := range gameState.AvailableNobles {
		if validator.CheckNobleVisit(playerState, &noble) {
			// Noble visits player
			playerState.Nobles = append(playerState.Nobles, noble)
			currentPlayer.VictoryPoints += noble.VictoryPoints

			// Remove noble from available
			gameState.AvailableNobles = append(gameState.AvailableNobles[:i], gameState.AvailableNobles[i+1:]...)
			break // Only one noble per turn
		}
	}

	// Update states
	if err := e.stateRepo.UpdateGameState(ctx, gameState); err != nil {
		return err
	}

	if err := e.stateRepo.UpdatePlayerState(ctx, playerState); err != nil {
		return err
	}

	if err := e.gameRepo.UpdatePlayer(ctx, currentPlayer); err != nil {
		return err
	}

	// Check victory condition
	if validator.CheckVictoryCondition(players) {
		// Someone reached 15 points - end the game
		winner := validator.DetermineWinner(players, map[int64]*models.PlayerState{})

		// Get all player states for proper winner determination
		playerStates := make(map[int64]*models.PlayerState)
		for _, p := range players {
			state, err := e.stateRepo.GetPlayerState(ctx, p.ID)
			if err == nil {
				playerStates[p.UserID] = state
			}
		}
		winner = validator.DetermineWinner(players, playerStates)

		// Update game status
		game.Status = models.GameStatusCompleted
		game.WinnerID = &winner.UserID
		now := time.Now()
		game.CompletedAt = &now

		if err := e.gameRepo.Update(ctx, game); err != nil {
			return fmt.Errorf("failed to update game status: %w", err)
		}

		return nil
	}

	// Switch turn
	if err := e.switchTurn(ctx, game, players); err != nil {
		return err
	}

	return nil
}

// ReserveCard executes the reserve card action
func (e *GameEngine) ReserveCard(ctx context.Context, gameID, userID int64, cardID int64, tier int) error {
	validator := NewGameValidator()

	// Get game
	game, err := e.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return err
	}

	// Validate turn
	if err := validator.ValidateTurn(game, userID); err != nil {
		return err
	}

	// Get player
	players, err := e.gameRepo.GetPlayers(ctx, gameID)
	if err != nil {
		return err
	}

	var currentPlayer *models.GamePlayer
	for _, p := range players {
		if p.UserID == userID {
			currentPlayer = p
			break
		}
	}

	// Get states
	gameState, err := e.stateRepo.GetGameState(ctx, gameID)
	if err != nil {
		return err
	}

	playerState, err := e.stateRepo.GetPlayerState(ctx, currentPlayer.ID)
	if err != nil {
		return err
	}

	// Validate can reserve
	if err := validator.ValidateReserveCard(playerState); err != nil {
		return err
	}

	var card *models.DevelopmentCard

	if cardID > 0 {
		// Reserve visible card
		c, err := e.cardRepo.GetCardByID(ctx, cardID)
		if err != nil {
			return err
		}
		card = c
		e.removeAndReplaceCard(gameState, card)
	} else {
		// Reserve from deck (blind)
		card = e.drawCardFromDeck(gameState, tier)
		if card == nil {
			return fmt.Errorf("no cards left in tier %d", tier)
		}
	}

	// Add to reserved
	playerState.ReservedCards = append(playerState.ReservedCards, *card)

	// Give gold coin if available
	if gameState.AvailableGems["gold"] > 0 {
		gameState.AvailableGems["gold"]--
		playerState.Gems["gold"]++
	}

	// Update states
	if err := e.stateRepo.UpdateGameState(ctx, gameState); err != nil {
		return err
	}

	if err := e.stateRepo.UpdatePlayerState(ctx, playerState); err != nil {
		return err
	}

	// Switch turn
	if err := e.switchTurn(ctx, game, players); err != nil {
		return err
	}

	return nil
}

// Helper functions
func (e *GameEngine) switchTurn(ctx context.Context, game *models.Game, players []*models.GamePlayer) error {
	// Find current player index
	currentIndex := -1
	for i, p := range players {
		if game.CurrentTurnPlayerID != nil && p.UserID == *game.CurrentTurnPlayerID {
			currentIndex = i
			break
		}
	}

	// Next player
	nextIndex := (currentIndex + 1) % len(players)
	game.CurrentTurnPlayerID = &players[nextIndex].UserID
	game.TurnNumber++

	return e.gameRepo.Update(ctx, game)
}

func (e *GameEngine) removeAndReplaceCard(gameState *models.GameState, card *models.DevelopmentCard) {
	switch card.Tier {
	case 1:
		newCards := []models.DevelopmentCard{}
		for _, c := range gameState.VisibleCardsTier1 {
			if c.ID != card.ID {
				newCards = append(newCards, c)
			}
		}
		gameState.VisibleCardsTier1 = newCards
		// Replace with new card from deck
		if len(gameState.DeckTier1) > 0 {
			gameState.VisibleCardsTier1 = append(gameState.VisibleCardsTier1, gameState.DeckTier1[0])
			gameState.DeckTier1 = gameState.DeckTier1[1:]
			gameState.DeckTier1Count--
		}
	case 2:
		newCards := []models.DevelopmentCard{}
		for _, c := range gameState.VisibleCardsTier2 {
			if c.ID != card.ID {
				newCards = append(newCards, c)
			}
		}
		gameState.VisibleCardsTier2 = newCards
		if len(gameState.DeckTier2) > 0 {
			gameState.VisibleCardsTier2 = append(gameState.VisibleCardsTier2, gameState.DeckTier2[0])
			gameState.DeckTier2 = gameState.DeckTier2[1:]
			gameState.DeckTier2Count--
		}
	case 3:
		newCards := []models.DevelopmentCard{}
		for _, c := range gameState.VisibleCardsTier3 {
			if c.ID != card.ID {
				newCards = append(newCards, c)
			}
		}
		gameState.VisibleCardsTier3 = newCards
		if len(gameState.DeckTier3) > 0 {
			gameState.VisibleCardsTier3 = append(gameState.VisibleCardsTier3, gameState.DeckTier3[0])
			gameState.DeckTier3 = gameState.DeckTier3[1:]
			gameState.DeckTier3Count--
		}
	}
}

func (e *GameEngine) drawCardFromDeck(gameState *models.GameState, tier int) *models.DevelopmentCard {
	switch tier {
	case 1:
		if len(gameState.DeckTier1) > 0 {
			card := gameState.DeckTier1[0]
			gameState.DeckTier1 = gameState.DeckTier1[1:]
			gameState.DeckTier1Count--
			return &card
		}
	case 2:
		if len(gameState.DeckTier2) > 0 {
			card := gameState.DeckTier2[0]
			gameState.DeckTier2 = gameState.DeckTier2[1:]
			gameState.DeckTier2Count--
			return &card
		}
	case 3:
		if len(gameState.DeckTier3) > 0 {
			card := gameState.DeckTier3[0]
			gameState.DeckTier3 = gameState.DeckTier3[1:]
			gameState.DeckTier3Count--
			return &card
		}
	}
	return nil
}
