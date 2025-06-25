# マルチステージビルドでサイズ最適化
FROM golang:1.23-bullseye AS builder

# 作業ディレクトリ設定
WORKDIR /app

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y \
    gcc \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係ダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# Swaggerドキュメント生成
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# ビルド（静的リンク）
RUN CGO_ENABLED=1 go build -ldflags="-w -s -extldflags '-static'" -o main .

# 実行用の軽量イメージ
FROM alpine:latest

# 必要なパッケージインストール
RUN apk --no-cache add ca-certificates sqlite

# 作業ディレクトリ作成
WORKDIR /app

# バイナリをコピー
COPY --from=builder /app/main .

# データディレクトリを作成
RUN mkdir -p ./data

# 実行権限を付与
RUN chmod +x ./main

# ポート公開
EXPOSE 3001

# 実行
CMD ["./main"] 