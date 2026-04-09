# ADR 0001: Gateway Language Choice

## Status
Accepted

## Context
Gateway traffic is I/O-heavy and requires low-latency request handling, robust middleware composition, and operational simplicity.

## Decision
Use Go 1.22 with Fiber v2 for the backend gateway.

## Consequences
- Positive: high throughput, simple deployment footprint, easy concurrency model.
- Negative: less compile-time domain safety than Rust; mitigated by focused testing and strict boundary validation.
