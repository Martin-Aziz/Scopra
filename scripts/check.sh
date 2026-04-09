#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR/backend"
gofmt -w ./cmd ./src ./tests
go test ./...

cd "$ROOT_DIR/cli"
cargo fmt
cargo test

cd "$ROOT_DIR/frontend"
npm run lint
npm run test
npm run build

echo "All checks passed."
