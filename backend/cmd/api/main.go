package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/config"
	"github.com/messederdaniel/real-time-notifications/backend/internal/database"
	"github.com/messederdaniel/real-time-notifications/backend/internal/routes"
	internalwebsocket "github.com/messederdaniel/real-time-notifications/backend/internal/websocket"
)

func main() {
	cfg := config.LoadConfig()

	dbPool, err := database.NewPostgresConnection(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	hub := internalwebsocket.NewHub()
	go hub.Run()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "http://192.168.")
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PATCH",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	routes.RegisterRoutes(router, dbPool, cfg.JWTSecret, hub)

	log.Printf("Starting server on port %s", cfg.ServerPort)

	if err := router.Run("0.0.0.0:" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
