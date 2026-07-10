package models

import "time"

type Notification struct {
	ID          string     `json:"id"`
	SenderID    string     `json:"sender_id"`
	RecipientID string     `json:"recipient_id"`
	ParentID    *string    `json:"parent_id,omitempty"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	IsRead      bool       `json:"is_read"`
	IsCompleted bool       `json:"is_completed"`
	CreatedAt   time.Time  `json:"created_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type NotificationWithUsers struct {
	Notification
	SenderName     string `json:"sender_name"`
	SenderEmail    string `json:"sender_email"`
	RecipientName  string `json:"recipient_name"`
	RecipientEmail string `json:"recipient_email"`
}

type CreateNotificationRequest struct {
	SenderID       string  `json:"sender_id" binding:"required"`
	RecipientID    string  `json:"recipient_id"`
	RecipientEmail string  `json:"recipient_email" binding:"required,email"`
	ParentID       *string `json:"parent_id,omitempty"`
	Title          string  `json:"title" binding:"required"`
	Message        string  `json:"message" binding:"required"`
}
