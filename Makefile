# ============================================================
# MUSEUM Development Makefile
# ============================================================
#
# 開発・テスト・コード生成・Lintなどの入口を
# このMakefileに統一する。
#
# 人間が実行しても、AIが実行しても、
# 同じコマンドで同じ結果になることを目的とする。
#
# 主なコマンド:
#
#   make setup
#   make dev
#   make generate
#   make test
#   make lint
#   make check
#
# ============================================================


.PHONY: \
	setup \
	dev \
	down \
	generate \
	backend-generate \
	generated-check \
	format \
	backend-format \
	test \
	backend-test \
	lint \
	backend-lint \
	frontend-lint \
	design-check \
	check


# ============================================================
# Setup
# ============================================================

# 初回セットアップ。
#
# 以下をまとめて実行する。
#
# - PostgreSQL起動
# - Backend依存関係取得
# - Frontend依存関係取得
# - OpenAPIコード生成
setup:
	docker compose up -d
	cd backend && go mod download
	cd frontend && npm install
	$(MAKE) generate


# ============================================================
# Development
# ============================================================

# Docker Compose上の開発用サービスを起動する。
#
# 現時点ではPostgreSQLなどの
# 開発用Infrastructureを起動する。
dev:
	docker compose up -d


# Docker Compose上のサービスを停止する。
down:
	docker compose down


# ============================================================
# Code Generation
# ============================================================

# すべての自動生成コードを更新する。
#
# 今後、
#
# - OpenAPI
# - mock
# - DB query
#
# などのコード生成が増えた場合も
# このコマンドを入口にする。
generate: backend-generate


# Backendのコード生成。
#
# 現在はOpenAPIから、
# Request / Response / Router / Server Interface
# などを生成する。
backend-generate:
	cd backend && go generate ./...


# OpenAPIなどの定義と生成コードが
# 一致しているか確認する。
#
# 例:
#
# openapi.yamlを変更
# ↓
# server.gen.goを更新し忘れる
# ↓
# git diffが発生
# ↓
# このコマンドが失敗する
#
# CIでも利用する想定。
generated-check:
	cd backend && go generate ./...
	git diff --exit-code


# ============================================================
# Format
# ============================================================

# プロジェクト全体のフォーマット。
format: backend-format


# Goコードをgofmtで整形する。
backend-format:
	cd backend && gofmt -w .


# ============================================================
# Test
# ============================================================

# プロジェクト全体のテスト。
#
# Frontendテストを導入した場合は
# ここに追加する。
test: backend-test


# Goの全テストを実行する。
backend-test:
	cd backend && go test ./...


# ============================================================
# Lint
# ============================================================

# プロジェクト全体のLint。
lint: backend-lint frontend-lint design-check


# Backendの静的チェック。
#
# go vetはGo標準の静的解析ツール。
backend-lint:
	cd backend && go vet ./...


# FrontendのLint。
frontend-lint:
	cd frontend && npm run lint


# Google Labs design.mdの仕様に沿って
# DESIGN.mdを検証する。
design-check:
	npx @google/design.md lint DESIGN.md


# ============================================================
# Full Check
# ============================================================

# コミット前に実行する標準チェック。
#
# 基本的には開発者もAIも、
# 修正完了後にこのコマンドを通す。
#
# 実行内容:
#
# 1. OpenAPIコード生成
# 2. Go format
# 3. Go test
# 4. Go vet
# 5. Frontend lint
# 6. DESIGN.md lint
#
# 将来的にCIもこのmake checkを呼ぶことで、
# ローカルとCIのチェック内容を統一する。
check:
	$(MAKE) generate
	$(MAKE) format
	$(MAKE) test
	$(MAKE) lint