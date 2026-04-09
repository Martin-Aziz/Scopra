package api

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martin-aziz/scopra/backend/src/middleware"
	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/repositories"
	"github.com/martin-aziz/scopra/backend/src/services"
	"github.com/martin-aziz/scopra/backend/src/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	auth       *services.AuthService
	toolCalls  *services.ToolCallService
	repository repositories.Repository
	db         *pgxpool.Pool
	cache      *redis.Client
	logger     utils.Logger
}

func NewHandler(
	auth *services.AuthService,
	toolCalls *services.ToolCallService,
	repository repositories.Repository,
	db *pgxpool.Pool,
	cache *redis.Client,
	logger utils.Logger,
) *Handler {
	return &Handler{
		auth:       auth,
		toolCalls:  toolCalls,
		repository: repository,
		db:         db,
		cache:      cache,
		logger:     logger,
	}
}

func RegisterRoutes(app *fiber.App, handler *Handler, tokenService *services.TokenService) {
	app.Get("/health", handler.liveness)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	apiV1 := app.Group("/api/v1")
	apiV1.Get("/health", handler.readiness)
	apiV1.Get("/dashboard/summary", handler.dashboardSummary)
	apiV1.Get("/audit/recent", handler.recentAuditEvents)

	authGroup := apiV1.Group("/auth", middleware.AuthRateLimit())
	authGroup.Post("/register", handler.register)
	authGroup.Post("/login", handler.login)
	authGroup.Post("/refresh", handler.refresh)

	protected := apiV1.Group("", middleware.APIRateLimit(), middleware.RequireAccessToken(tokenService))
	protected.Post("/tool-calls", handler.toolCall)
	protected.Post("/agents/:id/revoke", handler.revokeAgent)
	protected.Post("/connectors/:tool/connect", handler.connectTool)
}

type registerRequest struct {
	Email    string          `json:"email"`
	Password string          `json:"password"`
	Role     models.UserRole `json:"role"`
}

func (h *Handler) register(c *fiber.Ctx) error {
	var request registerRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, err := h.auth.Register(c.UserContext(), request.Email, request.Password, request.Role)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEmail) || errors.Is(err, services.ErrWeakPassword) || errors.Is(err, services.ErrInvalidRole) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "unable to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) login(c *fiber.Ctx) error {
	var request loginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, tokens, err := h.auth.Login(c.UserContext(), request.Email, request.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "authentication failed"})
	}

	setRefreshCookie(c, tokens.RefreshToken, int(tokens.RefreshTokenExpires))

	return c.JSON(fiber.Map{
		"accessToken":          tokens.AccessToken,
		"accessTokenExpiresIn": tokens.AccessTokenExpires,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) refresh(c *fiber.Ctx) error {
	refreshToken := strings.TrimSpace(c.Cookies("refresh_token"))
	if refreshToken == "" {
		var request refreshRequest
		if err := c.BodyParser(&request); err == nil {
			refreshToken = strings.TrimSpace(request.RefreshToken)
		}
	}

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}

	tokens, err := h.auth.Refresh(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	setRefreshCookie(c, tokens.RefreshToken, int(tokens.RefreshTokenExpires))
	return c.JSON(fiber.Map{
		"accessToken":          tokens.AccessToken,
		"accessTokenExpiresIn": tokens.AccessTokenExpires,
	})
}

func setRefreshCookie(c *fiber.Ctx, value string, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Path:     "/api/v1/auth/refresh",
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteStrictMode,
		MaxAge:   maxAge,
	})
}

func (h *Handler) toolCall(c *fiber.Ctx) error {
	var request models.ToolCallRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	correlationID := correlationIDFromContext(c)
	response, err := h.toolCalls.Handle(c.UserContext(), request, correlationID)
	if err != nil {
		h.logger.Warn("tool call failed", map[string]any{"error": err.Error(), "correlationId": correlationID})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}

func (h *Handler) revokeAgent(c *fiber.Ctx) error {
	if c.Locals("userRole") != string(models.RoleAdmin) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin role required"})
	}

	agentID := strings.TrimSpace(c.Params("id"))
	if _, err := uuid.Parse(agentID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid agent id"})
	}

	if err := h.repository.RevokeAgent(c.UserContext(), agentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to revoke agent"})
	}

	previousHash, _ := h.repository.GetLatestAuditHash(c.UserContext())
	_ = h.repository.InsertAuditEvent(c.UserContext(), models.AuditEvent{
		AgentID:       agentID,
		Tool:          "gateway",
		Action:        "agent.revoke",
		Status:        "success",
		Payload:       map[string]any{"revokedBy": c.Locals("userID")},
		CorrelationID: correlationIDFromContext(c),
	}, previousHash)

	return c.JSON(fiber.Map{"status": "revoked", "agentId": agentID})
}

type connectToolRequest struct {
	Scopes []string `json:"scopes"`
}

func (h *Handler) connectTool(c *fiber.Ctx) error {
	tool := strings.ToLower(strings.TrimSpace(c.Params("tool")))
	if !isSupportedTool(tool) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported tool"})
	}

	var request connectToolRequest
	if err := c.BodyParser(&request); err != nil {
		request.Scopes = []string{}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"tool":    tool,
		"status":  "connected",
		"scopes":  request.Scopes,
		"message": "connector registered",
	})
}

func isSupportedTool(tool string) bool {
	supported := map[string]struct{}{
		"github": {},
		"slack":  {},
		"jira":   {},
		"notion": {},
		"linear": {},
	}
	_, ok := supported[tool]
	return ok
}

func (h *Handler) liveness(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "service": "nexus-mcp-gateway"})
}

func (h *Handler) readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	health := fiber.Map{"status": "ready"}
	isReady := true

	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			health["database"] = "down"
			isReady = false
		} else {
			health["database"] = "up"
		}
	}

	if h.cache != nil {
		if err := h.cache.Ping(ctx).Err(); err != nil {
			health["cache"] = "down"
			isReady = false
		} else {
			health["cache"] = "up"
		}
	}

	if !isReady {
		health["status"] = "degraded"
		return c.Status(fiber.StatusServiceUnavailable).JSON(health)
	}

	return c.JSON(health)
}

func (h *Handler) dashboardSummary(c *fiber.Ctx) error {
	if h.db == nil {
		return c.JSON(fiber.Map{
			"requestsLast24h": 0,
			"activeAgents":    0,
			"connectedTools":  0,
			"blockedEvents":   0,
		})
	}

	var requestsLast24h, activeAgents, connectedTools, blockedEvents int

	_ = h.db.QueryRow(c.UserContext(), `SELECT COUNT(*) FROM audit_events WHERE created_at > NOW() - INTERVAL '24 hours'`).Scan(&requestsLast24h)
	_ = h.db.QueryRow(c.UserContext(), `SELECT COUNT(*) FROM agents WHERE status = 'active'`).Scan(&activeAgents)
	_ = h.db.QueryRow(c.UserContext(), `SELECT COUNT(*) FROM tool_connections WHERE status = 'connected'`).Scan(&connectedTools)
	_ = h.db.QueryRow(c.UserContext(), `SELECT COUNT(*) FROM audit_events WHERE status = 'failed' AND created_at > NOW() - INTERVAL '24 hours'`).Scan(&blockedEvents)

	return c.JSON(fiber.Map{
		"requestsLast24h": requestsLast24h,
		"activeAgents":    activeAgents,
		"connectedTools":  connectedTools,
		"blockedEvents":   blockedEvents,
	})
}

func (h *Handler) recentAuditEvents(c *fiber.Ctx) error {
	if h.db == nil {
		return c.JSON([]fiber.Map{})
	}

	rows, err := h.db.Query(c.UserContext(), `
		SELECT id, agent_id::text, tool_slug, action_name, status, correlation_id, created_at
		FROM audit_events
		ORDER BY id DESC
		LIMIT 50;
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch audit events"})
	}
	defer rows.Close()

	events := make([]fiber.Map, 0, 50)
	for rows.Next() {
		var (
			id            int64
			agentID       string
			tool          string
			action        string
			status        string
			correlationID sql.NullString
			createdAt     time.Time
		)
		if err := rows.Scan(&id, &agentID, &tool, &action, &status, &correlationID, &createdAt); err != nil {
			continue
		}
		events = append(events, fiber.Map{
			"id":            id,
			"agentId":       agentID,
			"tool":          tool,
			"action":        action,
			"status":        status,
			"correlationId": correlationID.String,
			"createdAt":     createdAt,
		})
	}

	return c.JSON(events)
}

func correlationIDFromContext(c *fiber.Ctx) string {
	correlationID, ok := c.Locals("correlationID").(string)
	if !ok || correlationID == "" {
		return "unknown"
	}
	return correlationID
}
