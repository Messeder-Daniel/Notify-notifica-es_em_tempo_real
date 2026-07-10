package services

import (
	"context"
	"strings"

	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) List(ctx context.Context, currentUserID string) ([]models.UserResponse, error) {
	currentUser, err := service.userRepository.FindByID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	if currentUser.Role != "admin" {
		return nil, ErrForbidden
	}

	return service.userRepository.List(ctx)
}

func (service *UserService) UpdateRole(ctx context.Context, currentUserID string, targetUserID string, role string) (*models.UserResponse, error) {
	role = strings.ToLower(strings.TrimSpace(role))

	if role != "admin" && role != "user" {
		return nil, ErrInvalidUserRole
	}

	currentUser, err := service.userRepository.FindByID(ctx, currentUserID)
	if err != nil {
		return nil, err
	}

	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	if currentUser.Role != "admin" {
		return nil, ErrForbidden
	}

	if currentUserID == targetUserID {
		return nil, ErrCannotChangeOwnRole
	}

	updatedUser, err := service.userRepository.UpdateRole(ctx, targetUserID, role)
	if err != nil {
		return nil, err
	}

	if updatedUser == nil {
		return nil, ErrUserNotFound
	}

	return updatedUser, nil
}
