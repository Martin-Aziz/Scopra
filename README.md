# NEXUS-MCP (Scopra)

## The Idea
NEXUS-MCP is an MCP gateway that lets AI agents call enterprise tools through a single secure control plane.
Instead of hand-building OAuth and audit plumbing per integration, teams get scoped access control, revocation, and traceable tool execution out of the box.
The MVP focuses on fast developer onboarding while preserving enterprise-grade security boundaries.

## Technical Architecture
The repository is organized as a monorepo with explicit boundaries:

- `backend/`: Go Fiber gateway with JWT auth, rate limiting, tool routing, approval queueing, and tamper-evident audit chain.
- `frontend/`: Next.js App Router dashboard for gateway health, KPI summaries, and recent audit visibility.
- `cli/`: Rust CLI for init, register, login, connect, deploy preflight, and agent revocation.
- `common/`: shared types/constants/schema contracts and OpenAPI source.
- `docs/`: architecture and ADR records.

### Data Flow
1. User authenticates via `/api/v1/auth/*`.
2. Agent submits a tool call to `/api/v1/tool-calls` with access token.
3. Gateway validates input and token, enforces rate limits, and routes to connector adapter.
4. Destructive calls can be queued for approval.
5. Every action is persisted to `audit_events` with cryptographic hash chaining.
6. Dashboard reads summary/audit projections for observability.

## Quick Start

### Prerequisites
- Go 1.22+
- Rust stable
- Node.js 20+
- Docker 24+

### Local bootstrap
```bash
cp .env.example .env
bash scripts/bootstrap.sh
```

### Start full stack
```bash
docker compose up --build
```

### Run gateway locally (without compose app container)
```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

### Run dashboard locally
```bash
cd frontend
cp .env.example .env
npm run dev
```

### Run CLI locally
```bash
cd cli
cargo run -- init --gateway-url http://localhost:8080
```

## Testing
```bash
# all checks
bash scripts/check.sh

# individual components
cd backend && go test ./...
cd cli && cargo test
cd frontend && npm run test
```

## MVP Features
- JWT authentication with access/refresh semantics and httpOnly refresh cookie on login.
- Auth and API rate limiting (`5/min` auth, `100/min` API).
- Connector routing for GitHub, Slack, Jira, Notion, and Linear via unified registry.
- Human-in-the-loop queue for destructive actions.
- Per-agent revocation endpoint.
- Cryptographically chained append-only audit events.
- Dashboard summary and recent audit stream.
- CLI workflow for initialization, auth, connector registration, deployment preflight, and revocation.

## Bonus Features Included
- Public, read-only dashboard summary and recent audit feed for first-run onboarding.
- `quick-connect` CLI command to generate a ready-to-use `nexus.yaml` manifest.
- Temporal approval window on destructive calls (15-minute expiry).

## Security Posture
- Input validation at API boundaries.
- JWT issuer and audience validation.
- bcrypt password hashing (cost factor 12).
- Security headers via Helmet.
- Generic external auth error messages.
- Structured JSON logging with correlation IDs.

## Documentation Index
- `DECISIONS.md`: assumptions, alternatives, and debt.
- `API.md`: endpoint and auth documentation.
- `PROMPTS.md`: prompt/process meta documentation.
- `docs/ARCHITECTURE.md`: system-level design.
- `docs/adr/`: architecture decision records.

## Known Limitations
- MVP connector execution is currently adapter-stubbed (no live SaaS OAuth exchange yet).
- SOC 2 controls are not yet fully implemented.
- Security audit step is informational in CI while dependency patching is in progress.