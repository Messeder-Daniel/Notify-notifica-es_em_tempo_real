package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/messederdaniel/real-time-notifications/backend/internal/models"
	"github.com/messederdaniel/real-time-notifications/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepository *repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (service *AuthService) Register(ctx context.Context, request models.RegisterRequest) (*models.LoginResponse, error) {
	name := strings.TrimSpace(request.Name)
	email := normalizeEmail(request.Email)
	password := strings.TrimSpace(request.Password)

	if name == "" {
		return nil, ErrUserNameRequired
	}

	if len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	existingUser, err := service.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, ErrUserEmailAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := service.userRepository.Create(ctx, name, email, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	token, err := service.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  models.NewUserResponse(*user),
	}, nil
}

func (service *AuthService) Login(ctx context.Context, request models.LoginRequest) (*models.LoginResponse, error) {
	request.Email = normalizeEmail(request.Email)

	user, err := service.ValidateCredentials(ctx, request)
	if err != nil {
		return nil, err
	}

	token, err := service.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		User:  models.NewUserResponse(*user),
	}, nil
}

func (service *AuthService) GetMe(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := service.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find authenticated user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	response := models.NewUserResponse(*user)
	return &response, nil
}

func (service *AuthService) UpdateProfile(ctx context.Context, userID string, request models.UpdateProfileRequest) (*models.UserResponse, error) {
	name := strings.TrimSpace(request.Name)
	email := normalizeEmail(request.Email)

	if name == "" {
		return nil, ErrUserNameRequired
	}

	currentUser, err := service.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find current user: %w", err)
	}

	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	if currentUser.Email != email {
		existingUser, err := service.userRepository.FindByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}

		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrUserEmailAlreadyExists
		}
	}

	updatedUser, err := service.userRepository.UpdateProfile(ctx, userID, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	if updatedUser == nil {
		return nil, ErrUserNotFound
	}

	response := models.NewUserResponse(*updatedUser)
	return &response, nil
}

func (service *AuthService) ChangePassword(ctx context.Context, userID string, request models.ChangePasswordRequest) error {
	newPassword := strings.TrimSpace(request.NewPassword)

	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	user, err := service.userRepository.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find current user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.CurrentPassword)); err != nil {
		return ErrCurrentPasswordInvalid
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := service.userRepository.UpdatePassword(ctx, userID, string(passwordHash)); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (service *AuthService) ValidateCredentials(ctx context.Context, request models.LoginRequest) (*models.User, error) {
	user, err := service.userRepository.FindByEmail(ctx, normalizeEmail(request.Email))
	if err != nil {
		return nil, fmt.Errorf("failed to validate credentials: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (service *AuthService) generateToken(user *models.User) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   now.Unix(),
		"exp":   expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(service.jwtSecret))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
