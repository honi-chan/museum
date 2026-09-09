# MUSEUM 開発用の共通コマンド。
#
# 人間が実行しても、AIエージェントが実行しても
# 同じ手順になるようにする。

.PHONY: setup dev test lint design-check backend-test frontend-lint

# 初回セットアップ。
setup:
	# PostgreSQLを起動する。
	docker compose up -d

	# Backend依存関係を取得する。
	cd backend && go mod download

	# Frontend依存関係を取得する。
	cd frontend && npm install

# 開発環境を起動する。
dev:
	docker compose up -d

# すべてのテストを実行する。
test: backend-test

# Backendのテスト。
backend-test:
	cd backend && go test ./...

# FrontendのLint。
frontend-lint:
	cd frontend && npm run lint

# DESIGN.mdを検証する。
#
# Google Labs design.md のルールに沿って、
# token参照・contrast・section構造などを確認する。
design-check:
	npx @google/design.md lint DESIGN.md

# コードとデザイン設計をまとめて確認する。
lint: frontend-lint design-check