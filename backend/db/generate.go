package db

//go:generate go tool sqlc generate -f ../sqlc.yaml

// このpackageはDBコード生成の入口。
//
// db/query/*.sql
// db/migrations/*.sql
//
// からsqlcがGoコードを生成する。
//
// 生成先:
//
// internal/infrastructure/database/sqlcgen/
//
// sqlcgen配下は自動生成コードなので
// 手動編集しない。
//
// DB変更時:
//
// 1. Migrationを追加
// 2. Query SQLを追加・変更
// 3. go generate ./...
// 4. Repositoryを実装
//
// の順番で進める。
