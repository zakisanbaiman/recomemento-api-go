.PHONY: build run dev test clean docs install

# 変数
BINARY_NAME=recomemento-api
MAIN_FILE=main.go
PORT=3001

# デフォルトターゲット
all: build

# 依存関係のインストール
install:
	go mod download
	go mod tidy

# Swaggerツールのインストール
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest

# ビルド
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

# 開発モード（ホットリロード）
dev:
	air

# 通常実行
run: build
	./$(BINARY_NAME)

# 直接実行
run-direct:
	go run $(MAIN_FILE)

# デバッグ実行（Delve）
debug:
	$(shell go env GOPATH)/bin/dlv debug $(MAIN_FILE) --listen=:2345 --headless --api-version=2 --accept-multiclient

# デバッグ実行（対話モード）
debug-interactive:
	$(shell go env GOPATH)/bin/dlv debug $(MAIN_FILE)

# テスト実行
test:
	go test -v ./...

# ユニットテストのみ実行
test-unit:
	go test -v ./handlers ./models ./dto

# 統合テストを含むテスト実行
test-integration:
	RUN_INTEGRATION_TESTS=1 go test -v ./...

# テストカバレッジ
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# 統合テスト含むカバレッジ
test-coverage-integration:
	RUN_INTEGRATION_TESTS=1 go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# 並列テスト実行
test-parallel:
	go test -v -parallel 4 ./...

# ベンチマークテスト
test-bench:
	go test -bench=. -benchmem ./...

# テスト詳細実行
test-verbose:
	go test -v -count=1 ./...

# 特定パッケージのテスト
test-handlers:
	go test -v ./handlers

test-models:
	go test -v ./models

test-integration-only:
	RUN_INTEGRATION_TESTS=1 go test -v -run "TestIntegration" ./...

# テスト結果をJUnit形式で出力（CI/CD用）
test-junit:
	go test -v ./... 2>&1 | go-junit-report > test-results.xml

# テストのwatch実行（airを使用）
test-watch:
	air -c .air.test.toml

# Swaggerドキュメント生成
docs:
	swag init

# コードフォーマット
fmt:
	go fmt ./...

# 静的解析
vet:
	go vet ./...

# リント（golangci-lintが必要）
lint:
	golangci-lint run

# クリーンアップ
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -rf docs/
	rm -rf tmp/

# データベースのリセット
reset-db:
	rm -f data/*.db
	rm -f prisma/dev.db*

# プロダクションビルド
build-prod:
	CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_FILE)

# 開発環境の初期セットアップ
setup: install install-tools docs
	@echo "開発環境のセットアップが完了しました"
	@echo "以下のコマンドで開発サーバーを起動できます:"
	@echo "make dev"

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo "  make setup        - 開発環境の初期セットアップ"
	@echo "  make install      - 依存関係のインストール"
	@echo "  make install-tools- 開発ツールのインストール"
	@echo "  make dev          - 開発モード（ホットリロード）"
	@echo "  make run          - アプリケーションの実行"
	@echo "  make run-direct   - 直接実行（ビルドなし）"
	@echo "  make debug        - デバッグ実行（リモートデバッグ）"
	@echo "  make debug-interactive - デバッグ実行（対話モード）"
	@echo "  make build        - ビルド"
	@echo "  make build-prod   - プロダクション用ビルド"
	@echo "  make test         - テスト実行"
	@echo "  make test-coverage- テストカバレッジ"
	@echo "  make docs         - Swaggerドキュメント生成"
	@echo "  make fmt          - コードフォーマット"
	@echo "  make vet          - 静的解析"
	@echo "  make lint         - リント実行"
	@echo "  make clean        - クリーンアップ"
	@echo "  make reset-db     - データベースリセット"
	@echo "  make help         - このヘルプ" 