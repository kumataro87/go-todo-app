# --- ビルドステージ ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 依存関係のキャッシュを先に取得（ソースより先にコピーすることでDockerレイヤーを活用）
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO無効、Linux向けの静的バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api .

# --- 実行ステージ（最小イメージ）---
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# バイナリとマイグレーションファイルをコピー
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/app/api"]
