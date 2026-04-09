package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martin-aziz/scopra/backend/src/models"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email string, passwordHash string, role models.UserRole) (models.User, error) {
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id::text, email, password_hash, role, is_active, created_at;
	`

	var user models.User
	if err := r.db.QueryRow(ctx, query, email, passwordHash, role).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt); err != nil {
		return models.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	query := `
		SELECT id::text, email, password_hash, role, is_active, created_at
		FROM users
		WHERE email = $1;
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) GetLatestAuditHash(ctx context.Context) (string, error) {
	query := `SELECT row_hash FROM audit_events ORDER BY id DESC LIMIT 1;`

	var hash string
	err := r.db.QueryRow(ctx, query).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("get latest audit hash: %w", err)
	}

	return hash, nil
}

func (r *PostgresRepository) InsertAuditEvent(ctx context.Context, event models.AuditEvent, previousHash string) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("marshal audit payload: %w", err)
	}

	query := `
		INSERT INTO audit_events (agent_id, tool_slug, action_name, status, payload, prev_hash, correlation_id)
		VALUES ($1::uuid, $2, $3, $4, $5::jsonb, $6, $7);
	`

	if _, err := r.db.Exec(ctx, query, event.AgentID, event.Tool, event.Action, event.Status, payload, previousHash, event.CorrelationID); err != nil {
		return fmt.Errorf("insert audit event: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CreateApprovalRequest(ctx context.Context, request models.ApprovalRequest) error {
	payload, err := json.Marshal(request.Payload)
	if err != nil {
		return fmt.Errorf("marshal approval payload: %w", err)
	}

	query := `
		INSERT INTO approval_requests (id, agent_id, tool_slug, action_name, payload, status, expires_at)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5::jsonb, $6, $7);
	`

	if _, err := r.db.Exec(ctx, query, request.ID, request.AgentID, request.Tool, request.Action, payload, request.Status, request.ExpiresAt); err != nil {
		return fmt.Errorf("create approval request: %w", err)
	}

	return nil
}

func (r *PostgresRepository) UpsertAgent(ctx context.Context, agentID string) error {
	query := `
		INSERT INTO agents (id, name, status)
		VALUES ($1::uuid, CONCAT('agent-', SUBSTRING($1::text, 1, 8)), 'active')
		ON CONFLICT (id)
		DO UPDATE SET status = 'active', updated_at = NOW();
	`

	if _, err := r.db.Exec(ctx, query, agentID); err != nil {
		return fmt.Errorf("upsert agent: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RevokeAgent(ctx context.Context, agentID string) error {
	query := `
		UPDATE agents
		SET status = 'revoked', updated_at = NOW()
		WHERE id = $1::uuid;
	`

	if _, err := r.db.Exec(ctx, query, agentID); err != nil {
		return fmt.Errorf("revoke agent: %w", err)
	}

	return nil
}

var ErrNotFound = errors.New("resource not found")
