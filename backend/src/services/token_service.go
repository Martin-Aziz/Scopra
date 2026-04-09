package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/martin-aziz/scopra/backend/src/models"
)

type TokenClaims struct {
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret     []byte
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenService(secret, issuer, audience string, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		secret:     []byte(secret),
		issuer:     issuer,
		audience:   audience,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *TokenService) IssueTokens(user models.User) (models.AuthTokens, error) {
	now := time.Now().UTC()
	access, err := s.signToken(user, "access", now.Add(s.accessTTL))
	if err != nil {
		return models.AuthTokens{}, err
	}

	refresh, err := s.signToken(user, "refresh", now.Add(s.refreshTTL))
	if err != nil {
		return models.AuthTokens{}, err
	}

	return models.AuthTokens{
		AccessToken:         access,
		AccessTokenExpires:  int64(s.accessTTL.Seconds()),
		RefreshToken:        refresh,
		RefreshTokenExpires: int64(s.refreshTTL.Seconds()),
	}, nil
}

func (s *TokenService) ParseAndValidate(tokenString string, expectedType string) (*TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(*TokenClaims)
	if !ok || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedType {
		return nil, ErrInvalidTokenType
	}

	if claims.Issuer != s.issuer {
		return nil, ErrInvalidToken
	}

	if !audienceContains(claims.Audience, s.audience) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *TokenService) signToken(user models.User, tokenType string, expiresAt time.Time) (string, error) {
	claims := TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      string(user.Role),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			Audience:  []string{s.audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", ErrTokenGeneration
	}
	return signed, nil
}

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrTokenGeneration  = errors.New("failed to generate token")
)

func audienceContains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
