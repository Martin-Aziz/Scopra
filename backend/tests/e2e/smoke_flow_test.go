package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/martin-aziz/scopra/backend/src/api"
	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/services"
	"github.com/martin-aziz/scopra/backend/src/utils"
)

type smokeRepository struct {
	users           map[string]models.User
	latestAuditHash string
}

func newSmokeRepository() *smokeRepository {
	return &smokeRepository{users: map[string]models.User{}}
}

func (m *smokeRepository) CreateUser(_ context.Context, email string, passwordHash string, role models.UserRole) (models.User, error) {
	user := models.User{ID: uuid.NewString(), Email: email, PasswordHash: passwordHash, Role: role, IsActive: true}
	m.users[email] = user
	return user, nil
}

func (m *smokeRepository) FindUserByEmail(_ context.Context, email string) (models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return user, nil
}

func (m *smokeRepository) GetLatestAuditHash(_ context.Context) (string, error) {
	return m.latestAuditHash, nil
}

func (m *smokeRepository) InsertAuditEvent(_ context.Context, _ models.AuditEvent, _ string) error {
	m.latestAuditHash = "hash"
	return nil
}

func (m *smokeRepository) CreateApprovalRequest(_ context.Context, _ models.ApprovalRequest) error {
	return nil
}

func (m *smokeRepository) UpsertAgent(_ context.Context, _ string) error {
	return nil
}

func (m *smokeRepository) RevokeAgent(_ context.Context, _ string) error {
	return nil
}

func TestRegisterLoginAndToolCallSmokeFlow(t *testing.T) {
	repository := newSmokeRepository()
	tokenService := services.NewTokenService("this-is-a-test-secret-with-at-least-32", "nexus", "clients", 15*time.Minute, 7*24*time.Hour)
	authService := services.NewAuthService(repository, tokenService, services.BcryptHasher{})
	toolService := services.NewToolCallService(repository, services.NewConnectorRegistry(), false)
	handler := api.NewHandler(authService, toolService, repository, nil, nil, utils.NewLogger("test"))

	app := fiber.New()
	api.RegisterRoutes(app, handler, tokenService)

	registerPayload := map[string]any{"email": "qa@example.com", "password": "very-strong-pass", "role": "admin"}
	registerBody, _ := json.Marshal(registerPayload)
	registerRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	registerRequest.Header.Set("Content-Type", "application/json")
	registerResponse, err := app.Test(registerRequest)
	if err != nil || registerResponse.StatusCode != http.StatusCreated {
		t.Fatalf("register flow failed: status=%d err=%v", registerResponse.StatusCode, err)
	}

	loginPayload := map[string]any{"email": "qa@example.com", "password": "very-strong-pass"}
	loginBody, _ := json.Marshal(loginPayload)
	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginResponse, err := app.Test(loginRequest)
	if err != nil || loginResponse.StatusCode != http.StatusOK {
		t.Fatalf("login flow failed: status=%d err=%v", loginResponse.StatusCode, err)
	}

	var loginResult struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(loginResponse.Body).Decode(&loginResult); err != nil {
		t.Fatalf("decode login response failed: %v", err)
	}

	toolPayload := map[string]any{
		"agentId":     uuid.NewString(),
		"tool":        "github",
		"action":      "repo.read",
		"payload":     map[string]any{"repository": "nexus-mcp"},
		"destructive": false,
	}
	toolBody, _ := json.Marshal(toolPayload)
	toolRequest := httptest.NewRequest(http.MethodPost, "/api/v1/tool-calls", bytes.NewReader(toolBody))
	toolRequest.Header.Set("Content-Type", "application/json")
	toolRequest.Header.Set("Authorization", "Bearer "+loginResult.AccessToken)
	toolResponse, err := app.Test(toolRequest)
	if err != nil || toolResponse.StatusCode != http.StatusOK {
		t.Fatalf("tool call flow failed: status=%d err=%v", toolResponse.StatusCode, err)
	}
}
