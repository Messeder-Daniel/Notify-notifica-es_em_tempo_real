package services

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrNotificationNotFound = errors.New("notification not found")
)
