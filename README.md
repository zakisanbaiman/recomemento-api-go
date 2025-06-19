# Recomemento API (Go版)

本の推薦システムのバックエンドAPI - Go言語で書き直された版

## 技術スタック

- **Go** 1.21+
- **Gin** - Webフレームワーク
- **GORM** - ORM
- **SQLite** - データベース
- **Swaggo** - Swagger ドキュメント生成

## 機能

- 本のCRUD操作
- 本の推薦機能（ジャンルと目的による）
- Swagger UIによるAPIドキュメント
- CORS対応
- ヘルスチェックエンドポイント

## 開発環境のセットアップ

### 前提条件

- Go 1.21以上がインストールされていること
- Git

### 1. リポジトリのクローンと依存関係のインストール

```bash
git clone [repository-url]
cd recomemento-api-go
go mod tidy
```

### 2. Swaggerドキュメントの生成

```bash
# swagツールのインストール（初回のみ）
go install github.com/swaggo/swag/cmd/swag@latest

# Swaggerドキュメントの生成
swag init
```

### 3. 環境変数の設定（オプション）

```bash
# .envファイルを作成（オプション）
echo 'DATABASE_URL="./data/books.db"' > .env
echo 'PORT="3001"' >> .env
```

## アプリケーションの実行

### 開発モード

```bash
# ホットリロード付きで実行（airツール使用の場合）
go install github.com/cosmtrek/air@latest
air

# または通常の実行
go run main.go
```

### 本番モード

```bash
# ビルド
go build -o recomemento-api main.go

# 実行
./recomemento-api
```

## APIドキュメント

アプリケーション起動後、以下のURLでSwagger UIにアクセスできます：

```
http://localhost:3001/api-docs/
```

OpenAPI JSON形式のドキュメントは以下から取得できます：

```
http://localhost:3001/api-json
```

## 利用可能なエンドポイント

### ヘルスチェック

- `GET /health` - APIの状態確認

### Books

- `POST /books` - 新しい本を作成
- `GET /books` - すべての本を取得
- `GET /books/:id` - 特定の本を取得
- `PATCH /books/:id` - 特定の本を更新
- `DELETE /books/:id` - 特定の本を削除
- `POST /books/recommend` - 本の推薦を取得

## プロジェクト構造

```
.
├── main.go              # アプリケーションのエントリーポイント
├── go.mod               # Goモジュール定義
├── models/              # データモデルとリポジトリ
│   └── book.go
├── handlers/            # HTTPハンドラー
│   └── book_handler.go
├── dto/                 # データ転送オブジェクト
│   └── book_dto.go
├── database/            # データベース設定とマイグレーション
│   └── database.go
├── docs/                # Swagger生成ファイル（自動生成）
└── data/                # SQLiteデータベースファイル
```

## 開発時のコマンド

```bash
# テストの実行
go test ./...

# コードフォーマット
go fmt ./...

# 静的解析
go vet ./...

# Swaggerドキュメントの再生成
swag init

# 依存関係の更新
go mod tidy
```

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `PORT` | `3001` | サーバーのポート番号 |
| `DATABASE_URL` | `./data/books.db` | SQLiteデータベースファイルのパス |

## TypeScript版からの主な変更点

1. **フレームワーク**: NestJS → Gin
2. **ORM**: Prisma → GORM
3. **データベース**: 同じSQLiteを使用（互換性維持）
4. **API仕様**: 元のAPIと完全互換
5. **Swagger**: NestJS Swagger → Swaggo
6. **パフォーマンス**: Goの高速性能とメモリ効率

## ライセンス

MIT

## Description

[Nest](https://github.com/nestjs/nest) framework TypeScript starter repository.

## Project setup

```bash
$ npm install
```

## Compile and run the project

```bash
# development
$ npm run start

# watch mode
$ npm run start:dev

# production mode
$ npm run start:prod
```

## Run tests

```bash
# unit tests
$ npm run test

# e2e tests
$ npm run test:e2e

# test coverage
$ npm run test:cov
```

