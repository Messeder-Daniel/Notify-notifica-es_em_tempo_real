package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Login(ctx *gin.Context) {
	var request models.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := handler.authService.Login(ctx.Request.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *AuthHandler) Me(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	userEmail, _ := ctx.Get("userEmail")

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"user_email": userEmail,
	})
}
