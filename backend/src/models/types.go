package models

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         UserRole
	IsActive     bool
	CreatedAt    time.Time
}

type AuthTokens struct {
	AccessToken         string `json:"accessToken"`
	AccessTokenExpires  int64  `json:"accessTokenExpiresIn"`
	RefreshToken        string `json:"refreshToken"`
	RefreshTokenExpires int64  `json:"refreshTokenExpiresIn"`
}

type ToolCallRequest struct {
	AgentID     string         `json:"agentId"`
	Tool        string         `json:"tool"`
	Action      string         `json:"action"`
	Payload     map[string]any `json:"payload"`
	Destructive bool           `json:"destructive"`
}

type ToolCallResponse struct {
	Status            string         `json:"status"`
	Result            map[string]any `json:"result,omitempty"`
	RequiresApproval  bool           `json:"requiresApproval"`
	ApprovalRequestID string         `json:"approvalRequestId,omitempty"`
	Message           string         `json:"message"`
}

type ApprovalRequest struct {
	ID        string
	AgentID   string
	Tool      string
	Action    string
	Payload   map[string]any
	Status    string
	ExpiresAt time.Time
}

type AuditEvent struct {
	AgentID       string
	Tool          string
	Action        string
	Status        string
	Payload       map[string]any
	CorrelationID string
}
