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
		RecipientEmail string `json:"recipient_email" binding:"required,email"`
		Title          string `json:"title" binding:"required"`
		Message        string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	senderID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notification, err := handler.notificationService.Create(ctx.Request.Context(), models.CreateNotificationRequest{
		SenderID:       senderID,
		RecipientEmail: request.RecipientEmail,
		Title:          request.Title,
		Message:        request.Message,
	})
	if err != nil {
		handler.handleNotificationError(ctx, err, "failed to create notification")
		return
	}

	handler.sendNotificationCreatedEvent(notification.RecipientID, notification)

	ctx.JSON(http.StatusCreated, notification)
}

func (handler *NotificationHandler) CreateReply(ctx *gin.Context) {
	var request struct {
		Title   string `json:"title" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	senderID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	parentID := ctx.Param("id")

	notification, err := handler.notificationService.CreateReply(
		ctx.Request.Context(),
		parentID,
		senderID,
		request.Title,
		request.Message,
	)
	if err != nil {
		handler.handleNotificationError(ctx, err, "failed to create reply")
		return
	}

	handler.sendNotificationCreatedEvent(notification.RecipientID, notification)

	ctx.JSON(http.StatusCreated, notification)
}

func (handler *NotificationHandler) FindReceivedByUserID(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notifications, err := handler.notificationService.ListReceivedByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (handler *NotificationHandler) FindSentByUserID(ctx *gin.Context) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notifications, err := handler.notificationService.ListSentByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find sent notifications"})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func (handler *NotificationHandler) MarkAsRead(ctx *gin.Context) {
	handler.updateNotificationStatus(ctx, "read")
}

func (handler *NotificationHandler) MarkAsUnread(ctx *gin.Context) {
	handler.updateNotificationStatus(ctx, "unread")
}

func (handler *NotificationHandler) MarkAsCompleted(ctx *gin.Context) {
	handler.updateNotificationStatus(ctx, "complete")
}

func (handler *NotificationHandler) Reopen(ctx *gin.Context) {
	handler.updateNotificationStatus(ctx, "reopen")
}

func (handler *NotificationHandler) updateNotificationStatus(ctx *gin.Context, action string) {
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	notificationID := ctx.Param("id")

	var (
		notification *models.NotificationWithUsers
		err          error
	)

	switch action {
	case "read":
		notification, err = handler.notificationService.MarkAsRead(ctx.Request.Context(), notificationID, userID)
	case "unread":
		notification, err = handler.notificationService.MarkAsUnread(ctx.Request.Context(), notificationID, userID)
	case "complete":
		notification, err = handler.notificationService.MarkAsCompleted(ctx.Request.Context(), notificationID, userID)
	case "reopen":
		notification, err = handler.notificationService.Reopen(ctx.Request.Context(), notificationID, userID)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification action"})
		return
	}

	if err != nil {
		handler.handleNotificationError(ctx, err, "failed to update notification")
		return
	}

	ctx.JSON(http.StatusOK, notification)
}

func (handler *NotificationHandler) handleNotificationError(ctx *gin.Context, err error, defaultMessage string) {
	switch {
	case errors.Is(err, services.ErrNotificationNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
	case errors.Is(err, services.ErrUserNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "recipient user not found"})
	case errors.Is(err, services.ErrNotificationTitleRequired),
		errors.Is(err, services.ErrNotificationMessageRequired):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": defaultMessage})
	}
}

func (handler *NotificationHandler) sendNotificationCreatedEvent(recipientID string, notification *models.NotificationWithUsers) {
	event := map[string]any{
		"type": "notification.created",
		"data": notification,
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal notification websocket event: %v", err)
		return
	}

	handler.hub.SendToUser(recipientID, eventData)
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
