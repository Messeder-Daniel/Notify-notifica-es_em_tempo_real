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
}

func NewNotificationService(notificationRepository *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (service *NotificationService) Create(ctx context.Context, request models.CreateNotificationRequest) (*models.Notification, error) {
	request.Title = strings.TrimSpace(request.Title)
	request.Message = strings.TrimSpace(request.Message)

	if request.Title == "" {
		return nil, ErrNotificationTitleRequired
	}

	if request.Message == "" {
		return nil, ErrNotificationMessageRequired
	}

	notification, err := service.notificationRepository.Create(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (service *NotificationService) FindByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	notifications, err := service.notificationRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications: %w", err)
	}

	return notifications, nil
}

func (service *NotificationService) MarkAsRead(ctx context.Context, notificationID string, userID string) (*models.Notification, error) {
	notification, err := service.notificationRepository.MarkAsRead(ctx, notificationID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	return notification, nil
}
