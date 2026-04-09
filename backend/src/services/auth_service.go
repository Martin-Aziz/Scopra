package services

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/repositories"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type BcryptHasher struct{}

func (h BcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (h BcryptHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type AuthService struct {
	repository   repositories.Repository
	tokenService *TokenService
	hasher       PasswordHasher
}

func NewAuthService(repository repositories.Repository, tokenService *TokenService, hasher PasswordHasher) *AuthService {
	return &AuthService{repository: repository, tokenService: tokenService, hasher: hasher}
}

func (s *AuthService) Register(ctx context.Context, email string, password string, role models.UserRole) (models.User, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return models.User{}, err
	}

	if err := validatePassword(password); err != nil {
		return models.User{}, err
	}

	if role == "" {
		role = models.RoleUser
	}

	if role != models.RoleAdmin && role != models.RoleUser {
		return models.User{}, ErrInvalidRole
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return models.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repository.CreateUser(ctx, normalizedEmail, hash, role)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (models.User, models.AuthTokens, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return models.User{}, models.AuthTokens{}, err
	}

	user, err := s.repository.FindUserByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return models.User{}, models.AuthTokens{}, ErrInvalidCredentials
		}
		return models.User{}, models.AuthTokens{}, err
	}

	if !user.IsActive {
		return models.User{}, models.AuthTokens{}, ErrInactiveUser
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return models.User{}, models.AuthTokens{}, ErrInvalidCredentials
	}

	tokens, err := s.tokenService.IssueTokens(user)
	if err != nil {
		return models.User{}, models.AuthTokens{}, err
	}

	return user, tokens, nil
}

func (s *AuthService) Refresh(refreshToken string) (models.AuthTokens, error) {
	claims, err := s.tokenService.ParseAndValidate(refreshToken, "refresh")
	if err != nil {
		return models.AuthTokens{}, ErrInvalidCredentials
	}

	user := models.User{
		ID:       claims.UserID,
		Email:    claims.Email,
		Role:     models.UserRole(claims.Role),
		IsActive: true,
	}

	tokens, err := s.tokenService.IssueTokens(user)
	if err != nil {
		return models.AuthTokens{}, err
	}

	return tokens, nil
}

func normalizeEmail(email string) (string, error) {
	candidate := strings.TrimSpace(strings.ToLower(email))
	if _, err := mail.ParseAddress(candidate); err != nil {
		return "", ErrInvalidEmail
	}
	return candidate, nil
}

func validatePassword(password string) error {
	if len(password) < 12 {
		return ErrWeakPassword
	}
	return nil
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrWeakPassword       = errors.New("password must be at least 12 characters")
	ErrInactiveUser       = errors.New("user is inactive")
	ErrInvalidRole        = errors.New("invalid user role")
)
