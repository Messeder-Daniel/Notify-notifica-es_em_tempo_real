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

func (repository *NotificationRepository) Create(ctx context.Context, request models.CreateNotificationRequest) (*models.NotificationWithUsers, error) {
	query := `
		WITH inserted AS (
			INSERT INTO notifications (
				sender_id,
				recipient_id,
				parent_id,
				title,
				message
			)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING *
		)
		SELECT
			n.id::text,
			n.sender_id::text,
			n.recipient_id::text,
			n.parent_id::text,
			n.title,
			n.message,
			n.is_read,
			n.is_completed,
			n.created_at,
			n.read_at,
			n.completed_at,
			sender.name,
			sender.email,
			recipient.name,
			recipient.email
		FROM inserted n
		INNER JOIN users sender ON sender.id = n.sender_id
		INNER JOIN users recipient ON recipient.id = n.recipient_id
	`

	return repository.scanNotification(
		repository.db.QueryRow(
			ctx,
			query,
			request.SenderID,
			request.RecipientID,
			request.ParentID,
			request.Title,
			request.Message,
		),
		"failed to create notification",
	)
}

func (repository *NotificationRepository) FindByID(ctx context.Context, notificationID string) (*models.NotificationWithUsers, error) {
	query := repository.baseSelectQuery() + `
		WHERE n.id = $1
	`

	notification, err := repository.scanNotification(
		repository.db.QueryRow(ctx, query, notificationID),
		"failed to find notification by id",
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return notification, nil
}

func (repository *NotificationRepository) ListReceivedByUserID(ctx context.Context, userID string) ([]models.NotificationWithUsers, error) {
	query := repository.baseSelectQuery() + `
		WHERE n.recipient_id = $1
		ORDER BY n.created_at DESC
	`

	return repository.queryNotifications(ctx, query, userID)
}

func (repository *NotificationRepository) ListSentByUserID(ctx context.Context, userID string) ([]models.NotificationWithUsers, error) {
	query := repository.baseSelectQuery() + `
		WHERE n.sender_id = $1
		ORDER BY n.created_at DESC
	`

	return repository.queryNotifications(ctx, query, userID)
}

func (repository *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	query := repository.updateSelectQuery(`
		is_read = TRUE,
		read_at = COALESCE(read_at, NOW())
	`)

	return repository.updateStatus(ctx, query, notificationID, userID, "failed to mark notification as read")
}

func (repository *NotificationRepository) MarkAsUnread(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	query := repository.updateSelectQuery(`
		is_read = FALSE,
		read_at = NULL
	`)

	return repository.updateStatus(ctx, query, notificationID, userID, "failed to mark notification as unread")
}

func (repository *NotificationRepository) MarkAsCompleted(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	query := repository.updateSelectQuery(`
		is_read = TRUE,
		read_at = COALESCE(read_at, NOW()),
		is_completed = TRUE,
		completed_at = COALESCE(completed_at, NOW())
	`)

	return repository.updateStatus(ctx, query, notificationID, userID, "failed to mark notification as completed")
}

func (repository *NotificationRepository) Reopen(ctx context.Context, notificationID string, userID string) (*models.NotificationWithUsers, error) {
	query := repository.updateSelectQuery(`
		is_completed = FALSE,
		completed_at = NULL
	`)

	return repository.updateStatus(ctx, query, notificationID, userID, "failed to reopen notification")
}

func (repository *NotificationRepository) updateStatus(
	ctx context.Context,
	query string,
	notificationID string,
	userID string,
	errorMessage string,
) (*models.NotificationWithUsers, error) {
	notification, err := repository.scanNotification(
		repository.db.QueryRow(ctx, query, notificationID, userID),
		errorMessage,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return notification, nil
}

func (repository *NotificationRepository) queryNotifications(ctx context.Context, query string, args ...any) ([]models.NotificationWithUsers, error) {
	rows, err := repository.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	notifications := make([]models.NotificationWithUsers, 0)

	for rows.Next() {
		notification, err := repository.scanNotification(rows, "failed to scan notification")
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, *notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notifications: %w", err)
	}

	return notifications, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (repository *NotificationRepository) scanNotification(row rowScanner, errorMessage string) (*models.NotificationWithUsers, error) {
	var notification models.NotificationWithUsers

	err := row.Scan(
		&notification.ID,
		&notification.SenderID,
		&notification.RecipientID,
		&notification.ParentID,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.IsCompleted,
		&notification.CreatedAt,
		&notification.ReadAt,
		&notification.CompletedAt,
		&notification.SenderName,
		&notification.SenderEmail,
		&notification.RecipientName,
		&notification.RecipientEmail,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errorMessage, err)
	}

	return &notification, nil
}

func (repository *NotificationRepository) baseSelectQuery() string {
	return `
		SELECT
			n.id::text,
			n.sender_id::text,
			n.recipient_id::text,
			n.parent_id::text,
			n.title,
			n.message,
			n.is_read,
			n.is_completed,
			n.created_at,
			n.read_at,
			n.completed_at,
			sender.name,
			sender.email,
			recipient.name,
			recipient.email
		FROM notifications n
		INNER JOIN users sender ON sender.id = n.sender_id
		INNER JOIN users recipient ON recipient.id = n.recipient_id
	`
}

func (repository *NotificationRepository) updateSelectQuery(setClause string) string {
	return fmt.Sprintf(`
		WITH updated AS (
			UPDATE notifications
			SET %s
			WHERE id = $1
			AND recipient_id = $2
			RETURNING *
		)
		SELECT
			n.id::text,
			n.sender_id::text,
			n.recipient_id::text,
			n.parent_id::text,
			n.title,
			n.message,
			n.is_read,
			n.is_completed,
			n.created_at,
			n.read_at,
			n.completed_at,
			sender.name,
			sender.email,
			recipient.name,
			recipient.email
		FROM updated n
		INNER JOIN users sender ON sender.id = n.sender_id
		INNER JOIN users recipient ON recipient.id = n.recipient_id
	`, setClause)
}
