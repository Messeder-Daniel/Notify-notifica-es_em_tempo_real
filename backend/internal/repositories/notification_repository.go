package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (repository *NotificationRepository) Create(ctx context.Context, request models.CreateNotificationRequest) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (user_id, title, message)
		VALUES ($1, $2, $3)
		RETURNING
			id::text,
			user_id::text,
			title,
			message,
			is_read,
			created_at,
			read_at
	`

	var notification models.Notification

	err := repository.db.QueryRow(
		ctx,
		query,
		request.UserID,
		request.Title,
		request.Message,
	).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
		&notification.ReadAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

func (repository *NotificationRepository) FindByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	query := `
		SELECT
			id::text,
			user_id::text,
			title,
			message,
			is_read,
			created_at,
			read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := repository.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find notifications by user id: %w", err)
	}
	defer rows.Close()

	notifications := make([]models.Notification, 0)

	for rows.Next() {
		var notification models.Notification

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Message,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.ReadAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notifications: %w", err)
	}

	return notifications, nil
}

func (repository *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) (*models.Notification, error) {
	query := `
		UPDATE notifications
		SET
			is_read = TRUE,
			read_at = NOW()
		WHERE id = $1
		AND user_id = $2
		RETURNING
			id::text,
			user_id::text,
			title,
			message,
			is_read,
			created_at,
			read_at
	`

	return repository.updateReadStatus(ctx, query, notificationID, userID, "failed to mark notification as read")
}

func (repository *NotificationRepository) MarkAsUnread(ctx context.Context, notificationID string, userID string) (*models.Notification, error) {
	query := `
		UPDATE notifications
		SET
			is_read = FALSE,
			read_at = NULL
		WHERE id = $1
		AND user_id = $2
		RETURNING
			id::text,
			user_id::text,
			title,
			message,
			is_read,
			created_at,
			read_at
	`

	return repository.updateReadStatus(ctx, query, notificationID, userID, "failed to mark notification as unread")
}

func (repository *NotificationRepository) updateReadStatus(
	ctx context.Context,
	query string,
	notificationID string,
	userID string,
	errorMessage string,
) (*models.Notification, error) {
	var notification models.Notification

	err := repository.db.QueryRow(ctx, query, notificationID, userID).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
		&notification.ReadAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: %w", errorMessage, err)
	}

	return &notification, nil
}
