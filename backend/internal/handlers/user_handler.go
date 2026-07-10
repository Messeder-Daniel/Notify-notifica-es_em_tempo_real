package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (handler *UserHandler) List(ctx *gin.Context) {
	currentUserID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	users, err := handler.userService.List(ctx.Request.Context(), currentUserID)
	if err != nil {
		handler.handleUserError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (handler *UserHandler) UpdateRole(ctx *gin.Context) {
	currentUserID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	var request struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	targetUserID := ctx.Param("id")

	user, err := handler.userService.UpdateRole(ctx.Request.Context(), currentUserID, targetUserID, request.Role)
	if err != nil {
		handler.handleUserError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) handleUserError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrForbidden):
		ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
	case errors.Is(err, services.ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, services.ErrInvalidUserRole):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "role must be admin or user"})
	case errors.Is(err, services.ErrCannotChangeOwnRole):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "you cannot change your own role"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user request"})
	}
}
