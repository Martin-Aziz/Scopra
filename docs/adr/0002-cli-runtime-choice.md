# ADR 0002: CLI Runtime Choice

## Status
Accepted

## Context
The CLI needs portable binaries, strong typed command parsing, and stable local configuration handling.

## Decision
Use Rust with clap + reqwest for the CLI runtime.

## Consequences
- Positive: predictable binaries, safe parsing, explicit error typing.
- Negative: additional toolchain for contributors; mitigated by bootstrap scripts and CI setup.
