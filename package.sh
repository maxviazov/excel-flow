#!/bin/bash

set -e

ARCHIVE_NAME="excel-flow-$(date +%Y%m%d-%H%M%S).tar.gz"

echo "📦 Создание архива проекта..."

tar -czf "$ARCHIVE_NAME" \
  --exclude='bin' \
  --exclude='*.exe' \
  --exclude='*.dll' \
  --exclude='*.so' \
  --exclude='*.dylib' \
  --exclude='flow' \
  --exclude='server' \
  --exclude='*.test' \
  --exclude='out' \
  --exclude='outputs' \
  --exclude='uploads' \
  --exclude='logs' \
  --exclude='.idea' \
  --exclude='.vscode' \
  --exclude='*.swp' \
  --exclude='*.swo' \
  --exclude='.DS_Store' \
  --exclude='Thumbs.db' \
  --exclude='terraform/.terraform' \
  --exclude='terraform/*.tfstate' \
  --exclude='terraform/*.tfstate.backup' \
  --exclude='terraform/.terraform.lock.hcl' \
  --exclude='terraform/terraform.tfvars' \
  --exclude='*.tmp' \
  --exclude='*.log' \
  --exclude='.git' \
  .

echo "✅ Архив создан: $ARCHIVE_NAME"
echo ""
echo "📋 Для распаковки на другой машине:"
echo "   tar -xzf $ARCHIVE_NAME"
echo "   cd excel-flow"
echo "   go mod download"
echo "   go run cmd/server/main.go"
