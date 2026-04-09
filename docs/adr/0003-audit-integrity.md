# ADR 0003: Audit Integrity Strategy

## Status
Accepted

## Context
Agent actions must be traceable and tamper-evident for governance and incident response.

## Decision
Persist audit events with chained hashes (`prev_hash` + generated `row_hash`) in PostgreSQL.

## Consequences
- Positive: low-friction tamper evidence and deterministic forensic trail.
- Negative: hash-chain continuity must be preserved during migrations and bulk backfills.
