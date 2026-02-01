#!/bin/bash

set -e

echo "==> Installing dependencies..."
go mod download

echo "==> Building project..."
go build -o server cmd/server/main.go

echo "==> Build completed successfully!"
