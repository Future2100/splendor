package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"splendor-backend/internal/domain/models"
	"splendor-backend/internal/repository/postgres"
)

var (
	ErrGameNotFound      = errors.New("game not found")
	ErrGameFull          = errors.New("game is full")
	ErrGameStarted       = errors.New("game has already started")
	ErrNotGameCreator    = errors.New("only the game creator can start the game")
	ErrNotEnoughPlayers  = errors.New("not enough players to start game")
	ErrAlreadyInGame     = errors.New("you are already in this game")
)

type GameService struct {
	gameRepo *postgres.GameRepository
	userRepo *postgres.UserRepository
	engine   GameEngine
}

type GameEngine interface {
	InitializeGame(ctx context.Context, gameID int64) error
	GetGameState(ctx context.Context, gameID int64) (*models.FullGameState, error)
}

func NewGameService(gameRepo *postgres.GameRepository, userRepo *postgres.UserRepository, engine GameEngine) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		userRepo: userRepo,
		engine:   engine,
	}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(ctx context.Context, userID int64, numPlayers int) (*models.CreateGameResponse, error) {
	// Generate unique room code
	roomCode, err := s.gameRepo.GenerateRoomCode(ctx)
	if err != nil {
		return nil, err
	}

	// Create game
	game := &models.Game{
		RoomCode:   roomCode,
		Status:     models.GameStatusWaiting,
		NumPlayers: numPlayers,
		CreatedBy:  userID,
	}

	if err := s.gameRepo.Create(ctx, game); err != nil {
		return nil, err
	}

	// Add creator as first player
	gamePlayer := &models.GamePlayer{
		GameID:         game.ID,
		UserID:         userID,
		PlayerPosition: 0,
		VictoryPoints:  0,
		IsActive:       true,
	}

	if err := s.gameRepo.AddPlayer(ctx, gamePlayer); err != nil {
		return nil, err
	}

	// Get full game details
	game, err = s.GetGameByID(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	return &models.CreateGameResponse{
		Game:     game,
		RoomCode: roomCode,
	}, nil
}

// JoinGame allows a user to join an existing game
func (s *GameService) JoinGame(ctx context.Context, userID int64, roomCode string) (*models.Game, error) {
	// Get game
	game, err := s.gameRepo.GetByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Check if game has started
	if game.Status != models.GameStatusWaiting {
		return nil, ErrGameStarted
	}

	// Check if user is already in game
	inGame, err := s.gameRepo.IsPlayerInGame(ctx, game.ID, userID)
	if err != nil {
		return nil, err
	}
	if inGame {
		return nil, ErrAlreadyInGame
	}

	// Check if game is full
	playerCount, err := s.gameRepo.GetPlayerCount(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	if playerCount >= game.NumPlayers {
		return nil, ErrGameFull
	}

	// Add player
	gamePlayer := &models.GamePlayer{
		GameID:         game.ID,
		UserID:         userID,
		PlayerPosition: playerCount,
		VictoryPoints:  0,
		IsActive:       true,
	}

	if err := s.gameRepo.AddPlayer(ctx, gamePlayer); err != nil {
		return nil, err
	}

	// Return updated game
	return s.GetGameByID(ctx, game.ID)
}

// LeaveGame removes a player from a game
func (s *GameService) LeaveGame(ctx context.Context, gameID, userID int64) error {
	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return ErrGameNotFound
	}

	// Can only leave if game is waiting
	if game.Status != models.GameStatusWaiting {
		return errors.New("cannot leave a game that has started")
	}

	return s.gameRepo.RemovePlayer(ctx, gameID, userID)
}

// StartGame starts a game
func (s *GameService) StartGame(ctx context.Context, gameID, userID int64) (*models.Game, error) {
	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Check if user is creator
	if game.CreatedBy != userID {
		return nil, ErrNotGameCreator
	}

	// Check if game already started
	if game.Status != models.GameStatusWaiting {
		return nil, ErrGameStarted
	}

	// Check minimum players (2)
	playerCount, err := s.gameRepo.GetPlayerCount(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	if playerCount < 2 {
		return nil, ErrNotEnoughPlayers
	}

	// Get players
	players, err := s.gameRepo.GetPlayers(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	// Update game status
	now := time.Now()
	game.Status = models.GameStatusInProgress
	game.StartedAt = &now
	game.CurrentTurnPlayerID = &players[0].UserID // First player starts

	if err := s.gameRepo.Update(ctx, game); err != nil {
		return nil, err
	}

	// Initialize game state
	if err := s.engine.InitializeGame(ctx, game.ID); err != nil {
		return nil, fmt.Errorf("failed to initialize game: %w", err)
	}

	// Get updated game with players
	return s.GetGameByID(ctx, game.ID)
}

// GetGameByID retrieves a game by ID with all players
func (s *GameService) GetGameByID(ctx context.Context, gameID int64) (*models.Game, error) {
	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// Get players
	players, err := s.gameRepo.GetPlayers(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	game.Players = players

	// Get creator
	creator, err := s.userRepo.GetByID(ctx, game.CreatedBy)
	if err == nil {
		game.Creator = creator
	}

	return game, nil
}

// GetGameByRoomCode retrieves a game by room code
func (s *GameService) GetGameByRoomCode(ctx context.Context, roomCode string) (*models.Game, error) {
	game, err := s.gameRepo.GetByRoomCode(ctx, roomCode)
	if err != nil {
		return nil, err
	}

	// Get players
	players, err := s.gameRepo.GetPlayers(ctx, game.ID)
	if err != nil {
		return nil, err
	}
	game.Players = players

	return game, nil
}

// ListGames lists all games with optional status filter
func (s *GameService) ListGames(ctx context.Context, status *models.GameStatus, limit, offset int) (*models.GameListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	games, total, err := s.gameRepo.List(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	// Populate players for each game
	for _, game := range games {
		players, err := s.gameRepo.GetPlayers(ctx, game.ID)
		if err == nil {
			game.Players = players
		}

		// Get creator
		creator, err := s.userRepo.GetByID(ctx, game.CreatedBy)
		if err == nil {
			game.Creator = creator
		}
	}

	return &models.GameListResponse{
		Games: games,
		Total: total,
	}, nil
}
