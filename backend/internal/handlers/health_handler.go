package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (handler *HealthHandler) Check(ctx *gin.Context) {
	pingCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := handler.db.Ping(pingCtx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "error",
			"message":  "API is running, but database is unavailable",
			"database": "disconnected",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "API is running",
		"database": "connected",
	})
}
