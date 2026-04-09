# NEXUS-MCP API Documentation

## Authentication
- Public endpoints: `/health`, `/api/v1/health`, `/api/v1/dashboard/summary`, `/api/v1/audit/recent`, `/api/v1/auth/*`
- Protected endpoints require `Authorization: Bearer <access_token>`.
- Access token lifetime: 15 minutes.
- Refresh token lifetime: 7 days (issued and rotated via `/api/v1/auth/refresh`, also set in httpOnly cookie during login).

## Endpoints

### POST /api/v1/auth/register
Creates a user account.

Request body:
```json
{
  "email": "admin@example.com",
  "password": "very-strong-password",
  "role": "admin"
}
```

Responses:
- 201 Created
- 400 Validation error
- 409 Conflict

### POST /api/v1/auth/login
Authenticates credentials and returns access token.

Request body:
```json
{
  "email": "admin@example.com",
  "password": "very-strong-password"
}
```

Responses:
- 200 OK
- 401 Unauthorized

### POST /api/v1/auth/refresh
Rotates refresh token and issues new access token.

Sources refresh token from cookie or body:
```json
{
  "refreshToken": "optional-if-cookie-present"
}
```

Responses:
- 200 OK
- 401 Unauthorized

### POST /api/v1/tool-calls
Executes or queues a tool action.

Request body:
```json
{
  "agentId": "f3d52fab-5913-4228-8fc4-109b2e30bb2d",
  "tool": "github",
  "action": "repo.read",
  "payload": {
    "repository": "nexus-mcp"
  },
  "destructive": false
}
```

Responses:
- 200 OK: executed or pending approval
- 400 Validation error
- 401 Unauthorized

### POST /api/v1/agents/{id}/revoke
Revokes an agent identity.

Authorization:
- Requires admin role.

Responses:
- 200 OK
- 403 Forbidden
- 400 Validation error

### POST /api/v1/connectors/{tool}/connect
Registers connector scopes.

Request body:
```json
{
  "scopes": ["repo.read", "issues.write"]
}
```

Responses:
- 201 Created
- 400 Unsupported connector
- 401 Unauthorized

### GET /api/v1/dashboard/summary
Returns aggregate metrics.

Response shape:
```json
{
  "requestsLast24h": 0,
  "activeAgents": 0,
  "connectedTools": 0,
  "blockedEvents": 0
}
```

### GET /api/v1/audit/recent
Returns latest 50 audit events.

### GET /health
Liveness endpoint.

### GET /api/v1/health
Readiness endpoint (DB + cache).

### GET /metrics
Prometheus metrics endpoint.

## Security Controls
- Auth routes rate-limited to 5 requests/minute.
- API routes rate-limited to 100 requests/minute.
- Passwords hashed with bcrypt (cost 12).
- Structured logs with correlation IDs.
- Generic auth errors to avoid account enumeration.
