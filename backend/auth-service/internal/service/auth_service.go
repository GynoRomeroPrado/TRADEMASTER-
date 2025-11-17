package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/trademaster/backend/auth-service/internal/models"
	"github.com/trademaster/backend/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(email, password, role, firstName, lastName, company, trade string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
	ValidateToken(tokenString string) (*models.User, error)
	GenerateToken(user *models.User) (string, error)
	GetUserByID(id string) (*models.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Signup(email, password, role, firstName, lastName, company, trade string) (*models.User, string, error) {
	existingUser, _ := s.repo.GetByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	profile := &models.UserProfile{
		ID:        uuid.New(),
		UserID:    user.ID,
		FirstName: firstName,
		LastName:  lastName,
		Company:   company,
		Trade:     trade,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateProfile(profile); err != nil {
		return nil, "", err
	}

	user.Profile = profile
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	profile, _ := s.repo.GetProfileByUserID(user.ID)
	user.Profile = profile

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	user.PasswordHash = ""
	return user, token, nil
}

func (s *authService) ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	profile, _ := s.repo.GetProfileByUserID(user.ID)
	user.Profile = profile
	user.PasswordHash = ""

	return user, nil
}

func (s *authService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) GetUserByID(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	profile, _ := s.repo.GetProfileByUserID(user.ID)
	user.Profile = profile
	user.PasswordHash = ""

	return user, nil
}
