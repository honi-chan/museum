# ============================================================
# MUSEUM Development Makefile
# ============================================================
#
# 開発者やAIエージェントが、
# 毎回長いコマンドを手入力しなくて済むように
# 開発・生成・テスト・DB操作を共通化する。
#
# 基本的には、
#
#   make backend-setup
#   make generate
#   make check
#
# を使えばよい。
#
# ============================================================

.PHONY: \
	setup \
	dev \
	generate \
	backend-generate \
	test \
	backend-test \
	lint \
	frontend-lint \
	design-check \
	check \
	generated-check \
	db-up \
	db-down \
	db-migrate-up \
	db-migrate-down \
	db-migrate-version \
	backend-setup


# ============================================================
# Setup
# ============================================================

# プロジェクト全体の初回セットアップ。
#
# PostgreSQL起動
# Backend依存関係取得
# Frontend依存関係取得
# コード生成
# Migration
#
# までまとめて実行する。
setup:
	docker compose up -d db
	cd backend && go mod download
	cd frontend && npm install
	cd backend && go generate ./...
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://museum:museum@db:5432/museum?sslmode=disable" \
		up


# ============================================================
# Development
# ============================================================

# 開発用サービスを起動する。
#
# 現時点ではPostgreSQLを起動する。
dev:
	docker compose up -d db


# ============================================================
# Code Generation
# ============================================================

# プロジェクト内の自動生成コードを更新する。
#
# 現在:
#
# OpenAPI
#   ↓
# oapi-codegen
#   ↓
# HTTP Server / Request / Response
#
# SQL
#   ↓
# sqlc
#   ↓
# PostgreSQL Query Code
#
# 将来コード生成が増えても、
# 開発者は make generate だけ実行すればよい。
generate: backend-generate frontend-generate


# Backendのコード生成。
#
# backend配下の //go:generate をすべて実行する。
backend-generate:
	cd backend && go generate ./...


# ============================================================
# Test
# ============================================================

# プロジェクト全体のテスト。
#
# 現時点ではBackendのみ。
# Frontendテスト導入後はここへ追加する。
test: backend-test


# Backendの全テストを実行する。
backend-test:
	cd backend && go test ./...


# ============================================================
# Lint
# ============================================================

# プロジェクト全体のLint。
lint: frontend-lint design-check


# FrontendのLint。
frontend-lint:
	cd frontend && npm run lint


# DESIGN.mdをGoogle design.md仕様で検証する。
#
# デザイントークン、
# section構造、
# token参照などをチェックする。
design-check:
	npx @google/design.md lint DESIGN.md


# ============================================================
# Check
# ============================================================

# コミット前の標準チェック。
#
# 1. 自動生成
# 2. Go Format
# 3. Backend Test
# 4. Frontend Lint
# 5. DESIGN.md Check
#
# をまとめて実行する。
#
# 基本的にはコミット前にこれを実行する。
check:
	cd backend && go generate ./...
	cd backend && gofmt -w .
	cd backend && go test ./...
	cd frontend && npm run lint
	npx @google/design.md lint DESIGN.md


# ============================================================
# Generated Code Check
# ============================================================

# OpenAPIやSQLと、
# 自動生成コードが一致しているか確認する。
#
# 例えば、
#
# openapi.yamlを変更
# ↓
# server.gen.goを更新し忘れる
#
# といった事故を検出できる。
#
# CIでもこのTargetを使用できる。
generated-check:
	cd backend && go generate ./...
	git diff --exit-code


# ============================================================
# Database
# ============================================================

# PostgreSQLを起動する。
db-up:
	docker compose up -d db


# Docker Compose環境を停止する。
#
# DB Volumeは削除しないため、
# データは保持される。
db-down:
	docker compose down


# 未適用Migrationをすべて適用する。
db-migrate-up:
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://museum:museum@db:5432/museum?sslmode=disable" \
		up


# Migrationを1つ戻す。
#
# 開発中に直前のMigrationを戻したい場合に使用する。
db-migrate-down:
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://museum:museum@db:5432/museum?sslmode=disable" \
		down 1


# 現在適用されているMigration Versionを確認する。
db-migrate-version:
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://museum:museum@db:5432/museum?sslmode=disable" \
		version


# ============================================================
# Backend Setup
# ============================================================

# Backend開発環境をまとめて準備する。
#
# PostgreSQL起動
# ↓
# Migration
# ↓
# OpenAPI / sqlcコード生成
# ↓
# Backend Test
#
# 新しいPCや、
# Backend環境を作り直した際は
# 基本的にこれを実行すればよい。
backend-setup:
	docker compose up -d db
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://museum:museum@db:5432/museum?sslmode=disable" \
		up
	cd backend && go generate ./...
	cd backend && go test ./...
.PHONY: frontend-generate
frontend-generate:
	cd frontend && npm run generate:api
