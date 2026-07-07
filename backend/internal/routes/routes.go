package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/messederdaniel/real-time-notifications/backend/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, db *pgxpool.Pool) {
	healthHandler := handlers.NewHealthHandler(db)

	router.GET("/health", healthHandler.Check)
}
