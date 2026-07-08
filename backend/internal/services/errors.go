package services

import "errors"

var (
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrNotificationNotFound        = errors.New("notification not found")
	ErrNotificationTitleRequired   = errors.New("notification title is required")
	ErrNotificationMessageRequired = errors.New("notification message is required")
)
