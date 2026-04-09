package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/services"
)

type memoryRepository struct {
	users           map[string]models.User
	latestAuditHash string
	approvals       []models.ApprovalRequest
	audits          []models.AuditEvent
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{users: map[string]models.User{}}
}

func (m *memoryRepository) CreateUser(_ context.Context, email string, passwordHash string, role models.UserRole) (models.User, error) {
	user := models.User{ID: uuid.NewString(), Email: email, PasswordHash: passwordHash, Role: role, IsActive: true}
	m.users[email] = user
	return user, nil
}

func (m *memoryRepository) FindUserByEmail(_ context.Context, email string) (models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return models.User{}, errors.New("not found")
	}
	return user, nil
}

func (m *memoryRepository) GetLatestAuditHash(_ context.Context) (string, error) {
	return m.latestAuditHash, nil
}

func (m *memoryRepository) InsertAuditEvent(_ context.Context, event models.AuditEvent, _ string) error {
	m.audits = append(m.audits, event)
	m.latestAuditHash = "hash"
	return nil
}

func (m *memoryRepository) CreateApprovalRequest(_ context.Context, request models.ApprovalRequest) error {
	m.approvals = append(m.approvals, request)
	return nil
}

func (m *memoryRepository) UpsertAgent(_ context.Context, _ string) error {
	return nil
}

func (m *memoryRepository) RevokeAgent(_ context.Context, _ string) error {
	return nil
}

func TestDestructiveToolCallQueuesApproval(t *testing.T) {
	repo := newMemoryRepository()
	registry := services.NewConnectorRegistry()
	service := services.NewToolCallService(repo, registry, true)

	request := models.ToolCallRequest{
		AgentID:     uuid.NewString(),
		Tool:        "github",
		Action:      "issues.write",
		Payload:     map[string]any{"issue": "test"},
		Destructive: true,
	}

	response, err := service.Handle(context.Background(), request, "corr-1")
	if err != nil {
		t.Fatalf("expected tool call to queue approval: %v", err)
	}

	if !response.RequiresApproval {
		t.Fatal("expected response to require approval")
	}

	if len(repo.approvals) != 1 {
		t.Fatalf("expected 1 approval request, got %d", len(repo.approvals))
	}
}

func TestNonDestructiveToolCallExecutes(t *testing.T) {
	repo := newMemoryRepository()
	registry := services.NewConnectorRegistry()
	service := services.NewToolCallService(repo, registry, false)

	request := models.ToolCallRequest{
		AgentID:     uuid.NewString(),
		Tool:        "slack",
		Action:      "channels.read",
		Payload:     map[string]any{"channel": "ops"},
		Destructive: false,
	}

	response, err := service.Handle(context.Background(), request, "corr-2")
	if err != nil {
		t.Fatalf("expected tool call to execute: %v", err)
	}

	if response.Status != "success" {
		t.Fatalf("expected success response, got %s", response.Status)
	}

	if len(repo.audits) == 0 {
		t.Fatal("expected audit event to be recorded")
	}
}
