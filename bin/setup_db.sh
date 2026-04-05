#!/bin/bash

# スクリプトが保存されているディレクトリを取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# .env ファイルから環境変数を読み込む
if [ -f "$PROJECT_ROOT/.env" ]; then
  set -a
  source "$PROJECT_ROOT/.env"
  set +e
else
  echo "Error: .env file not found at $PROJECT_ROOT/.env"
  exit 1
fi

# 環境変数の確認
if [ -z "$PROJECT_NAME" ] || [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_DB" ]; then
  echo "Error: Required environment variables not set in .env"
  echo "Required: PROJECT_NAME, POSTGRES_USER, POSTGRES_DB"
  exit 1
fi

echo "Setting up database: $POSTGRES_DB..."
echo "Container: ${PROJECT_NAME}_postgres"

# マイグレーションファイルの実行
docker exec "${PROJECT_NAME}_postgres" psql \
  -U "${POSTGRES_USER}" \
  -d "${POSTGRES_DB}" \
  -f /var/lib/postgresql/migration/create_db.sql

if [ $? -eq 0 ]; then
  echo "✓ Database setup completed successfully"
else
  echo "✗ Database setup failed"
  exit 1
fi