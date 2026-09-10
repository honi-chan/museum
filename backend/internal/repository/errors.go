package repository

import "errors"

// ErrNotFound はRepository共通の
// 「対象データが存在しない」エラー。
//
// PostgreSQLのpgx.ErrNoRowsなど
// Infrastructure固有のエラーを
// UseCaseへ漏らさないために使用する。
var ErrNotFound = errors.New(
	"repository entity not found",
)
