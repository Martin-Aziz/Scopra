package repositories

import (
	"context"

	"github.com/martin-aziz/scopra/backend/src/models"
)

type Repository interface {
	CreateUser(ctx context.Context, email string, passwordHash string, role models.UserRole) (models.User, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
	GetLatestAuditHash(ctx context.Context) (string, error)
	InsertAuditEvent(ctx context.Context, event models.AuditEvent, previousHash string) error
	CreateApprovalRequest(ctx context.Context, request models.ApprovalRequest) error
	UpsertAgent(ctx context.Context, agentID string) error
	RevokeAgent(ctx context.Context, agentID string) error
}
