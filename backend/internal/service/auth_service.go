package service

import (
	"context"
	"errors"
	"fmt"

	"splendor-backend/internal/domain/models"
	"splendor-backend/internal/repository/postgres"
	"splendor-backend/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type AuthService struct {
	userRepo          *postgres.UserRepository
	jwtSecret         string
	jwtAccessExpiry   int64
	jwtRefreshExpiry  int64
}

func NewAuthService(userRepo *postgres.UserRepository, jwtSecret string, accessExpiry, refreshExpiry int64) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtAccessExpiry:  accessExpiry,
		jwtRefreshExpiry: refreshExpiry,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrEmailAlreadyExists
	}

	// Check if username already exists
	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, ErrUsernameAlreadyExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := jwt.GenerateTokenPair(
		user.ID,
		user.Username,
		s.jwtSecret,
		s.jwtAccessExpiry,
		s.jwtRefreshExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := jwt.GenerateTokenPair(
		user.ID,
		user.Username,
		s.jwtSecret,
		s.jwtAccessExpiry,
		s.jwtRefreshExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}

// RefreshToken generates new tokens from a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := jwt.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	tokenPair, err := jwt.GenerateTokenPair(
		user.ID,
		user.Username,
		s.jwtSecret,
		s.jwtAccessExpiry,
		s.jwtRefreshExpiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         user,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
