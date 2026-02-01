package api

import (
	"splendor-backend/internal/api/handlers"
	"splendor-backend/internal/api/middleware"
	"splendor-backend/internal/config"
	"splendor-backend/internal/gamelogic"
	"splendor-backend/internal/repository/postgres"
	"splendor-backend/internal/service"
	"splendor-backend/pkg/database"
	"splendor-backend/pkg/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *database.DB, hub *websocket.Hub, cfg *config.Config) {
	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	gameRepo := postgres.NewGameRepository(db)
	cardRepo := postgres.NewCardRepository(db)
	stateRepo := postgres.NewStateRepository(db)
	statsRepo := postgres.NewStatsRepository(db)

	// Initialize game engine
	gameEngine := gamelogic.NewGameEngine(gameRepo, cardRepo, stateRepo)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)
	gameService := service.NewGameService(gameRepo, userRepo, gameEngine)
	statsService := service.NewStatsService(statsRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	gameHandler := handlers.NewGameHandler(gameService)
	wsHandler := handlers.NewWebSocketHandler(hub, cfg.JWTSecret)
	stateHandler := handlers.NewStateHandler(gameEngine)
	gameplayHandler := handlers.NewGameplayHandler(gameEngine, hub)
	statsHandler := handlers.NewStatsHandler(statsService)
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Splendor API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetCurrentUser)
		}

		// Game routes
		games := v1.Group("/games")
		{
			games.GET("", gameHandler.ListGames)
			games.POST("", middleware.AuthMiddleware(cfg.JWTSecret), gameHandler.CreateGame)
			games.POST("/join", middleware.AuthMiddleware(cfg.JWTSecret), gameHandler.JoinGame)
			games.GET("/:id", gameHandler.GetGame)
			games.GET("/:id/state", stateHandler.GetGameState)
			games.POST("/:id/leave", middleware.AuthMiddleware(cfg.JWTSecret), gameHandler.LeaveGame)
			games.POST("/:id/start", middleware.AuthMiddleware(cfg.JWTSecret), gameHandler.StartGame)

			// Gameplay actions
			games.POST("/:id/take-gems", middleware.AuthMiddleware(cfg.JWTSecret), gameplayHandler.TakeGems)
			games.POST("/:id/purchase-card", middleware.AuthMiddleware(cfg.JWTSecret), gameplayHandler.PurchaseCard)
			games.POST("/:id/reserve-card", middleware.AuthMiddleware(cfg.JWTSecret), gameplayHandler.ReserveCard)
		}

		// WebSocket route
		v1.GET("/ws/games/:id", wsHandler.HandleConnection)

		// Stats routes
		stats := v1.Group("/stats")
		{
			stats.GET("/users/:id", statsHandler.GetUserStats)
			stats.GET("/leaderboard", statsHandler.GetLeaderboard)
		}
	}
}
