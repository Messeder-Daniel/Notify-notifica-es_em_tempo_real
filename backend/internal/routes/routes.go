package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/messederdaniel/real-time-notifications/backend/internal/handlers"
	"github.com/messederdaniel/real-time-notifications/backend/internal/middlewares"
	"github.com/messederdaniel/real-time-notifications/backend/internal/repositories"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
	internalwebsocket "github.com/messederdaniel/real-time-notifications/backend/internal/websocket"
)

func RegisterRoutes(router *gin.Engine, db *pgxpool.Pool, jwtSecret string, hub *internalwebsocket.Hub) {
	healthHandler := handlers.NewHealthHandler(db)

	userRepository := repositories.NewUserRepository(db)
	notificationRepository := repositories.NewNotificationRepository(db)

	authService := services.NewAuthService(userRepository, jwtSecret)
	notificationService := services.NewNotificationService(notificationRepository)

	authHandler := handlers.NewAuthHandler(authService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	webSocketHandler := handlers.NewWebSocketHandler(hub, jwtSecret)

	router.GET("/health", healthHandler.Check)
	router.POST("/auth/login", authHandler.Login)
	router.GET("/ws", webSocketHandler.Connect)

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(jwtSecret))

	protected.GET("/auth/me", authHandler.Me)
	protected.GET("/notifications", notificationHandler.FindByUserID)
	protected.POST("/notifications", notificationHandler.Create)
	protected.PATCH("/notifications/:id/read", notificationHandler.MarkAsRead)
}
