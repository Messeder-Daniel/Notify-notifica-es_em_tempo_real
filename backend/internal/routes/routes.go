package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/messederdaniel/real-time-notifications/backend/internal/handlers"
	"github.com/messederdaniel/real-time-notifications/backend/internal/middlewares"
	"github.com/messederdaniel/real-time-notifications/backend/internal/repositories"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
)

func RegisterRoutes(router *gin.Engine, db *pgxpool.Pool, jwtSecret string) {
	healthHandler := handlers.NewHealthHandler(db)

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router.GET("/health", healthHandler.Check)
	router.POST("/auth/login", authHandler.Login)

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))
	protected.GET("/auth/me", authHandler.Me)
}
