package gamelogic

import (
	"errors"

	"splendor-backend/internal/domain/models"
)

var (
	ErrNotYourTurn        = errors.New("it's not your turn")
	ErrInvalidGemCount    = errors.New("invalid gem count")
	ErrNotEnoughGems      = errors.New("not enough gems available")
	ErrTooManyGems        = errors.New("too many gems (max 10)")
	ErrInvalidCardTier    = errors.New("invalid card tier")
	ErrCardNotAvailable   = errors.New("card not available")
	ErrCannotAffordCard   = errors.New("cannot afford this card")
	ErrTooManyReserved    = errors.New("too many reserved cards (max 3)")
)

type GameValidator struct{}

func NewGameValidator() *GameValidator {
	return &GameValidator{}
}

// ValidateTurn checks if it's the player's turn
func (v *GameValidator) ValidateTurn(game *models.Game, userID int64) error {
	if game.CurrentTurnPlayerID == nil || *game.CurrentTurnPlayerID != userID {
		return ErrNotYourTurn
	}
	return nil
}

// ValidateTakeGems validates taking gems action
func (v *GameValidator) ValidateTakeGems(gameState *models.GameState, playerState *models.PlayerState, gems map[string]int) error {
	totalTaking := 0
	differentColors := 0
	sameColorCount := 0
	var sameColorType string

	for gemType, count := range gems {
		if count < 0 {
			return ErrInvalidGemCount
		}
		if count > 0 {
			if gemType == "gold" {
				return errors.New("cannot take gold coins directly")
			}
			totalTaking += count
			differentColors++
			if count > 1 {
				sameColorCount = count
				sameColorType = gemType
			}
		}
	}

	// Rule: Take 3 different colors OR 2 of the same color
	if differentColors == 3 && sameColorCount == 0 && totalTaking == 3 {
		// Taking 3 different colors - valid
		for gemType, count := range gems {
			if count > 0 && gameState.AvailableGems[gemType] < count {
				return ErrNotEnoughGems
			}
		}
	} else if differentColors == 1 && sameColorCount == 2 && totalTaking == 2 {
		// Taking 2 of the same color - valid only if at least 4 available
		if gameState.AvailableGems[sameColorType] < 4 {
			return errors.New("need at least 4 gems to take 2 of same color")
		}
	} else {
		return ErrInvalidGemCount
	}

	// Check 10 gem limit
	currentTotal := 0
	for _, count := range playerState.Gems {
		currentTotal += count
	}
	if currentTotal+totalTaking > 10 {
		return ErrTooManyGems
	}

	return nil
}

// ValidatePurchaseCard validates purchasing a card
func (v *GameValidator) ValidatePurchaseCard(gameState *models.GameState, playerState *models.PlayerState, card *models.DevelopmentCard) error {
	// Calculate total gold needed
	totalGoldNeeded := 0

	for gemType, cost := range card.Cost {
		if cost > 0 {
			permanent := playerState.PermanentGems[gemType]
			owned := playerState.Gems[gemType]
			available := permanent + owned

			if available < cost {
				// Need gold to cover the difference
				totalGoldNeeded += cost - available
			}
		}
	}

	// Check if player has enough gold
	if totalGoldNeeded > playerState.Gems["gold"] {
		return ErrCannotAffordCard
	}

	return nil
}

// ValidateReserveCard validates reserving a card
func (v *GameValidator) ValidateReserveCard(playerState *models.PlayerState) error {
	if len(playerState.ReservedCards) >= 3 {
		return ErrTooManyReserved
	}
	return nil
}

// CalculateCost calculates the actual cost including gold usage
func (v *GameValidator) CalculateCost(card *models.DevelopmentCard, playerState *models.PlayerState) map[string]int {
	actualCost := make(map[string]int)
	goldNeeded := 0

	for gemType, cost := range card.Cost {
		if cost > 0 {
			permanent := playerState.PermanentGems[gemType]
			remaining := cost - permanent
			if remaining > 0 {
				gems := playerState.Gems[gemType]
				if gems >= remaining {
					actualCost[gemType] = remaining
				} else {
					actualCost[gemType] = gems
					goldNeeded += remaining - gems
				}
			}
		}
	}

	if goldNeeded > 0 {
		actualCost["gold"] = goldNeeded
	}

	return actualCost
}

// CheckNobleVisit checks if a noble should visit the player
func (v *GameValidator) CheckNobleVisit(playerState *models.PlayerState, noble *models.Noble) bool {
	for gemType, required := range noble.Required {
		if playerState.PermanentGems[gemType] < required {
			return false
		}
	}
	return true
}

// CheckVictoryCondition checks if game should end (15+ points)
func (v *GameValidator) CheckVictoryCondition(players []*models.GamePlayer) bool {
	for _, player := range players {
		if player.VictoryPoints >= 15 {
			return true
		}
	}
	return false
}

// DetermineWinner determines the winner based on points and tiebreaker
func (v *GameValidator) DetermineWinner(players []*models.GamePlayer, playerStates map[int64]*models.PlayerState) *models.GamePlayer {
	var winner *models.GamePlayer
	maxPoints := -1
	minCards := 999999

	for _, player := range players {
		if player.VictoryPoints > maxPoints {
			maxPoints = player.VictoryPoints
			winner = player
			if state, ok := playerStates[player.UserID]; ok {
				minCards = len(state.PurchasedCards)
			}
		} else if player.VictoryPoints == maxPoints {
			// Tiebreaker: fewer cards wins
			if state, ok := playerStates[player.UserID]; ok {
				if len(state.PurchasedCards) < minCards {
					winner = player
					minCards = len(state.PurchasedCards)
				}
			}
		}
	}

	return winner
}
