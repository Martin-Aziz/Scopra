# Architecture Overview

## Bounded Contexts
- Identity and Access: registration, login, token issuance, refresh, role checks.
- Routing and Connectors: tool call validation, connector dispatch, result handling.
- Approval and Governance: destructive call queueing, expiry windows, admin revocation.
- Audit and Observability: chained audit persistence, metrics exposure, health checks.

## Layering
### Backend
- `api/`: HTTP transport and boundary parsing.
- `services/`: business orchestration and policy decisions.
- `repositories/`: persistence and data access concerns.
- `middleware/`: cross-cutting auth/rate-limit/correlation/metrics concerns.
- `models/`: domain data contracts.

### Frontend
- `app/`: route-level pages.
- `components/`: reusable UI building blocks.
- `services/`: API fetching and adapters.
- `utils/`: rendering and class composition helpers.

### CLI
- `commands/`: user-facing command workflows.
- `client/`: HTTP interface to gateway.
- `config/`: local machine configuration state.

## Data Storage
- PostgreSQL: users, agents, tool_connections, approval_requests, audit_events.
- Redis: readiness and future queue/cache support.

## Security Boundaries
- Gateway enforces JWT issuer/audience and role checks.
- Mutations require bearer token and route-level limits.
- Audit trail receives entries for critical security actions and tool-call outcomes.
