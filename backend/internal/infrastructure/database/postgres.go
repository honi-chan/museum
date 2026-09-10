package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgres はPostgreSQL connection poolを生成する。
//
// DB接続処理をmain.goから分離することで、
// main.goは依存関係の組み立てだけにする。
func NewPostgres(
	ctx context.Context,
	databaseURL string,
) (*pgxpool.Pool, error) {
	// DATABASE_URLからconnection poolを生成する。
	pool, err := pgxpool.New(
		ctx,
		databaseURL,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create postgres pool: %w",
			err,
		)
	}

	// Pool生成だけでは実際にDBへ接続しないため、
	// Pingして起動時に接続確認する。
	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"ping postgres: %w",
			err,
		)
	}

	return pool, nil
}
