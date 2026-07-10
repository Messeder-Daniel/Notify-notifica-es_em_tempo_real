package services

import "errors"

var (
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrUserNotFound                = errors.New("user not found")
	ErrUserNameRequired            = errors.New("user name is required")
	ErrUserEmailAlreadyExists      = errors.New("email already registered")
	ErrCurrentPasswordInvalid      = errors.New("current password is invalid")
	ErrPasswordTooShort            = errors.New("password must have at least 6 characters")
	ErrForbidden                   = errors.New("forbidden")
	ErrInvalidUserRole             = errors.New("invalid user role")
	ErrCannotChangeOwnRole         = errors.New("cannot change your own role")
	ErrNotificationNotFound        = errors.New("notification not found")
	ErrNotificationTitleRequired   = errors.New("notification title is required")
	ErrNotificationMessageRequired = errors.New("notification message is required")
)
