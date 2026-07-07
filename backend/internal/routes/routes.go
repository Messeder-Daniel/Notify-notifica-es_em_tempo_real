package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/handlers"
)

func RegisterRoutes(router *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	router.GET("/health", healthHandler.Check)
}
