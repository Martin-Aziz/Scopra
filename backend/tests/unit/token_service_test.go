package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/services"
)

func TestIssueAndValidateAccessToken(t *testing.T) {
	tokenService := services.NewTokenService("this-is-a-test-secret-with-at-least-32", "nexus", "clients", 15*time.Minute, 7*24*time.Hour)
	user := models.User{ID: uuid.NewString(), Email: "dev@example.com", Role: models.RoleUser}

	tokens, err := tokenService.IssueTokens(user)
	if err != nil {
		t.Fatalf("expected issue tokens to succeed: %v", err)
	}

	claims, err := tokenService.ParseAndValidate(tokens.AccessToken, "access")
	if err != nil {
		t.Fatalf("expected parse and validate to succeed: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected user id %s, got %s", user.ID, claims.UserID)
	}

	if claims.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, claims.Email)
	}
}

func TestRefreshTokenRejectedWhenAccessExpected(t *testing.T) {
	tokenService := services.NewTokenService("this-is-a-test-secret-with-at-least-32", "nexus", "clients", 15*time.Minute, 7*24*time.Hour)
	user := models.User{ID: uuid.NewString(), Email: "dev@example.com", Role: models.RoleAdmin}

	tokens, err := tokenService.IssueTokens(user)
	if err != nil {
		t.Fatalf("expected issue tokens to succeed: %v", err)
	}

	_, err = tokenService.ParseAndValidate(tokens.RefreshToken, "access")
	if err == nil {
		t.Fatal("expected refresh token validation to fail for access token expectation")
	}
}
