package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/trademaster/backend/auth-service/internal/models"
	"github.com/trademaster/backend/auth-service/internal/repository"
	"github.com/trademaster/backend/shared/config"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, password, role string) (*models.User, error)
	Login(email, password string) (string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetCurrentUser(userID uuid.UUID) (*models.User, error)
	Logout(userID uuid.UUID) error
}

type authService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
	jwtConfig   *config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, redisClient *redis.Client, jwtConfig *config.JWTConfig) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
		jwtConfig:   jwtConfig,
	}
}

func (s *authService) Register(email, password, role string) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		Subscription: models.Subscription{
			Plan:   "free",
			Status: "active",
		},
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	// Store refresh token in Redis
	ctx := context.Background()
	s.redisClient.Set(ctx, "refresh_token:"+user.ID.String(), refreshToken, 30*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	// Validate refresh token
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return "", err
	}

	// Check if refresh token exists in Redis
	ctx := context.Background()
	storedToken, err := s.redisClient.Get(ctx, "refresh_token:"+userID.String()).Result()
	if err != nil || storedToken != refreshToken {
		return "", errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})
}

func (s *authService) GetCurrentUser(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) Logout(userID uuid.UUID) error {
	ctx := context.Background()
	return s.redisClient.Del(ctx, "refresh_token:"+userID.String()).Err()
}

func (s *authService) generateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *authService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}
