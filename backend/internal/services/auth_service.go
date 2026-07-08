package services

import (
	"context"
	"fmt"
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

func (service *AuthService) Login(ctx context.Context, request models.LoginRequest) (*models.LoginResponse, error) {
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

func (service *AuthService) ValidateCredentials(ctx context.Context, request models.LoginRequest) (*models.User, error) {
	user, err := service.userRepository.FindByEmail(ctx, request.Email)
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
