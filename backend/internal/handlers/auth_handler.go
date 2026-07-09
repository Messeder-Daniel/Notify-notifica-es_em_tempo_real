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

func (handler *AuthHandler) Register(ctx *gin.Context) {
	var request models.RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := handler.authService.Register(ctx.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNameRequired),
			errors.Is(err, services.ErrPasswordTooShort):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case errors.Is(err, services.ErrUserEmailAlreadyExists):
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, response)
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
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	response, err := handler.authService.GetMe(ctx.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find authenticated user"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *AuthHandler) UpdateMe(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	var request models.UpdateProfileRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := handler.authService.UpdateProfile(ctx.Request.Context(), userID, request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNameRequired):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case errors.Is(err, services.ErrUserEmailAlreadyExists):
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		case errors.Is(err, services.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *AuthHandler) ChangePassword(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	var request models.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := handler.authService.ChangePassword(ctx.Request.Context(), userID, request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPasswordTooShort),
			errors.Is(err, services.ErrCurrentPasswordInvalid):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case errors.Is(err, services.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
