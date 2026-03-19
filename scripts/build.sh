#!/usr/bin/env bash
set -euo pipefail

target="${1:-all}"

build_assets() {
  echo "Building assets"
  npm run build
}

build_server() {
  echo "Building server"
  go build -o bin/server ./cmd/server
}

case "$target" in
  assets)
    build_assets
    ;;
  server)
    build_server
    ;;
  all|"")
    build_assets
    build_server
    ;;
  *)
    echo "Usage: $0 [assets|server|all]"
    exit 1
    ;;
esac