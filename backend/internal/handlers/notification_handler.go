package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (handler *NotificationHandler) Create(ctx *gin.Context) {
	var request struct {
		Title   string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notification, err := handler.notificationService.Create(ctx.Request.Context(), models.CreateNotificationRequest{
		UserID:  userID,
		Title:   request.Title,
		Message: request.Message,
	})
	if err != nil {
		if errors.Is(err, services.ErrNotificationTitleRequired) || errors.Is(err, services.ErrNotificationMessageRequired) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	ctx.JSON(http.StatusCreated, notification)
}

func (handler *NotificationHandler) FindByUserID(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notifications, err := handler.notificationService.FindByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (handler *NotificationHandler) MarkAsRead(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notificationID := ctx.Param("id")

	notification, err := handler.notificationService.MarkAsRead(ctx.Request.Context(), notificationID, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotificationNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark notification as read"})
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func getAuthenticatedUserID(ctx *gin.Context) (string, bool) {
	value, exists := ctx.Get("userID")
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}
