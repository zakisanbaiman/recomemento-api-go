# マルチステージビルドでサイズ最適化
FROM golang:1.23-alpine AS builder

# 作業ディレクトリ設定
WORKDIR /app

# 必要なパッケージをインストール
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 依存関係ファイルをコピー
COPY go.mod go.sum ./

# 依存関係ダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# Swaggerドキュメント生成
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# ビルド（最適化オプション付き）
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# 実行用の軽量イメージ
FROM alpine:latest

# 必要なパッケージインストール
RUN apk --no-cache add ca-certificates sqlite

# 作業ディレクトリ作成
WORKDIR /root/

# バイナリをコピー
COPY --from=builder /app/main .

# データディレクトリを作成
RUN mkdir -p ./data

# ポート公開
EXPOSE 3001

# 実行
CMD ["./main"] 