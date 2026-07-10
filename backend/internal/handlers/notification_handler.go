package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/services"
	internalwebsocket "github.com/messederdaniel/real-time-notifications/backend/internal/websocket"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
	hub                 *internalwebsocket.Hub
}

func NewNotificationHandler(notificationService *services.NotificationService, hub *internalwebsocket.Hub) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		hub:                 hub,
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

	handler.sendNotificationCreatedEvent(userID, notification)

	ctx.JSON(http.StatusCreated, notification)
}

func (handler *NotificationHandler) CreateForRecipientEmail(ctx *gin.Context) {
	var request struct {
		RecipientEmail string `json:"recipient_email" binding:"required,email"`
		Title          string `json:"title" binding:"required"`
		Message        string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if _, ok := getAuthenticatedUserID(ctx); !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notification, err := handler.notificationService.CreateForRecipientEmail(
		ctx.Request.Context(),
		request.RecipientEmail,
		request.Title,
		request.Message,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "recipient user not found"})
			return
		case errors.Is(err, services.ErrNotificationTitleRequired),
			errors.Is(err, services.ErrNotificationMessageRequired):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send notification"})
			return
		}
	}

	handler.sendNotificationCreatedEvent(notification.UserID, notification)

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
	handler.updateReadStatus(ctx, true)
}

func (handler *NotificationHandler) MarkAsUnread(ctx *gin.Context) {
	handler.updateReadStatus(ctx, false)
}

func (handler *NotificationHandler) updateReadStatus(ctx *gin.Context, shouldMarkAsRead bool) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notificationID := ctx.Param("id")

	var (
		notification *models.Notification
		err          error
	)

	if shouldMarkAsRead {
		notification, err = handler.notificationService.MarkAsRead(ctx.Request.Context(), notificationID, userID)
	} else {
		notification, err = handler.notificationService.MarkAsUnread(ctx.Request.Context(), notificationID, userID)
	}

	if err != nil {
		if errors.Is(err, services.ErrNotificationNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func (handler *NotificationHandler) sendNotificationCreatedEvent(userID string, notification *models.Notification) {
	event := map[string]any{
		"type": "notification.created",
		"data": notification,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal notification websocket event: %v", err)
		return
	}

	handler.hub.SendToUser(userID, eventData)
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
