#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR"

echo "[1/4] Starting infrastructure containers"
docker compose up -d db cache

echo "[2/4] Installing backend dependencies"
cd "$ROOT_DIR/backend"
go mod tidy

echo "[3/4] Installing CLI dependencies"
cd "$ROOT_DIR/cli"
cargo fetch

echo "[4/4] Installing frontend dependencies"
cd "$ROOT_DIR/frontend"
npm install

echo "Bootstrap complete."
