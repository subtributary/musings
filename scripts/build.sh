#!/usr/bin/env bash
set -e

echo "Building server"
go build -o bin/server ./cmd/server

