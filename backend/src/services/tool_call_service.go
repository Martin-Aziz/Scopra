package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/repositories"
)

type ToolCallService struct {
	repository   repositories.Repository
	connectors   *ConnectorRegistry
	approvalMode bool
}

func NewToolCallService(repository repositories.Repository, connectors *ConnectorRegistry, approvalMode bool) *ToolCallService {
	return &ToolCallService{repository: repository, connectors: connectors, approvalMode: approvalMode}
}

func (s *ToolCallService) Handle(ctx context.Context, request models.ToolCallRequest, correlationID string) (models.ToolCallResponse, error) {
	if err := validateToolCallRequest(request); err != nil {
		return models.ToolCallResponse{}, err
	}

	if err := s.repository.UpsertAgent(ctx, request.AgentID); err != nil {
		return models.ToolCallResponse{}, err
	}

	if request.Destructive && s.approvalMode {
		return s.queueApproval(ctx, request, correlationID)
	}

	result, err := s.connectors.Execute(ctx, strings.ToLower(request.Tool), request.Action, request.Payload)
	if err != nil {
		s.writeAuditEvent(ctx, request, correlationID, "failed")
		return models.ToolCallResponse{}, err
	}

	s.writeAuditEvent(ctx, request, correlationID, "success")

	return models.ToolCallResponse{
		Status:           "success",
		Result:           result,
		RequiresApproval: false,
		Message:          "tool call executed",
	}, nil
}

func (s *ToolCallService) queueApproval(ctx context.Context, request models.ToolCallRequest, correlationID string) (models.ToolCallResponse, error) {
	approvalID := uuid.NewString()
	approval := models.ApprovalRequest{
		ID:        approvalID,
		AgentID:   request.AgentID,
		Tool:      request.Tool,
		Action:    request.Action,
		Payload:   request.Payload,
		Status:    "pending",
		ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
	}

	if err := s.repository.CreateApprovalRequest(ctx, approval); err != nil {
		return models.ToolCallResponse{}, err
	}

	s.writeAuditEvent(ctx, request, correlationID, "pending_approval")

	return models.ToolCallResponse{
		Status:            "pending",
		RequiresApproval:  true,
		ApprovalRequestID: approvalID,
		Message:           "destructive action queued for approval",
	}, nil
}

func (s *ToolCallService) writeAuditEvent(ctx context.Context, request models.ToolCallRequest, correlationID string, status string) {
	previousHash, err := s.repository.GetLatestAuditHash(ctx)
	if err != nil {
		return
	}

	audit := models.AuditEvent{
		AgentID:       request.AgentID,
		Tool:          strings.ToLower(request.Tool),
		Action:        request.Action,
		Status:        status,
		Payload:       request.Payload,
		CorrelationID: correlationID,
	}

	_ = s.repository.InsertAuditEvent(ctx, audit, previousHash)
}

func validateToolCallRequest(request models.ToolCallRequest) error {
	if _, err := uuid.Parse(request.AgentID); err != nil {
		return ErrInvalidAgentID
	}

	if strings.TrimSpace(request.Tool) == "" {
		return ErrToolRequired
	}

	if strings.TrimSpace(request.Action) == "" {
		return ErrActionRequired
	}

	if request.Payload == nil {
		return ErrPayloadRequired
	}

	return nil
}

var (
	ErrInvalidAgentID  = errors.New("invalid agentId")
	ErrToolRequired    = errors.New("tool is required")
	ErrActionRequired  = errors.New("action is required")
	ErrPayloadRequired = errors.New("payload is required")
)
