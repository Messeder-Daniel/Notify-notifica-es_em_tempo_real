package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/routes"
)

func main() {
	router := gin.Default()

	routes.RegisterRoutes(router)

	log.Println("Starting server on port 8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
