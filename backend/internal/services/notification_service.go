package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/repositories"
)

type NotificationService struct {
	notificationRepository *repositories.NotificationRepository
	userRepository         *repositories.UserRepository
}

func NewNotificationService(
	notificationRepository *repositories.NotificationRepository,
	userRepository *repositories.UserRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
		userRepository:         userRepository,
	}
}

func (service *NotificationService) Create(ctx context.Context, request models.CreateNotificationRequest) (*models.NotificationWithUsers, error) {
	request.Title = strings.TrimSpace(request.Title)
	request.Message = strings.TrimSpace(request.Message)
	request.RecipientEmail = strings.ToLower(strings.TrimSpace(request.RecipientEmail))

	if request.Title == "" {
		return nil, ErrNotificationTitleRequired
	}

	if request.Message == "" {
		return nil, ErrNotificationMessageRequired
	}

	if request.RecipientID == "" {
		recipient, err := service.userRepository.FindByEmail(ctx, request.RecipientEmail)
		if err != nil {
			return nil, fmt.Errorf("failed to find recipient: %w", err)
		}

		if recipient == nil {
			return nil, ErrUserNotFound
		}

		request.RecipientID = recipient.ID
	}

	notification, err := service.notificationRepository.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (service *NotificationService) CreateReply(
	ctx context.Context,
	parentID string,
	senderID string,
	title string,
	message string,
) (*models.NotificationWithUsers, error) {
	parentNotification, err := service.notificationRepository.FindByID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find parent notification: %w", err)
	}

	if parentNotification == nil {
		return nil, ErrNotificationNotFound
	}

	recipientID := parentNotification.SenderID
	if senderID == parentNotification.SenderID {
		recipientID = parentNotification.RecipientID
	}

	return service.Create(ctx, models.CreateNotificationRequest{
		SenderID:    senderID,
		RecipientID: recipientID,
		ParentID:    &parentID,
		Title:       title,
		Message:     message,
	})
}

func (service *NotificationService) ListReceivedByUserID(ctx context.Context, userID string) ([]models.NotificationWithUsers, error) {
	notifications, err := service.notificationRepository.ListReceivedByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list received notifications: %w", err)
	}

	return notifications, nil
}

func (service *NotificationService) ListSentByUserID(ctx context.Context, userID string) ([]models.NotificationWithUsers, error) {
	notifications, err := service.notificationRepository.ListSentByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list sent notifications: %w", err)
	}

	return notifications, nil
}

func (service *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	notification, err := service.notificationRepository.MarkAsRead(ctx, notificationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}

func (service *NotificationService) MarkAsUnread(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	notification, err := service.notificationRepository.MarkAsUnread(ctx, notificationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as unread: %w", err)
	}

	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}

func (service *NotificationService) MarkAsCompleted(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	notification, err := service.notificationRepository.MarkAsCompleted(ctx, notificationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as completed: %w", err)
	}

	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}

func (service *NotificationService) Reopen(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	notification, err := service.notificationRepository.Reopen(ctx, notificationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen notification: %w", err)
	}

	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}
