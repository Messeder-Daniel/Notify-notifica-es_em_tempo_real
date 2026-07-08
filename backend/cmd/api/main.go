package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/config"
	"github.com/messederdaniel/real-time-notifications/backend/internal/database"
	"github.com/messederdaniel/real-time-notifications/backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	dbPool, err := database.NewPostgresConnection(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	routes.RegisterRoutes(router, dbPool, cfg.JWTSecret)

	log.Printf("Starting server on port %s", cfg.ServerPort)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
