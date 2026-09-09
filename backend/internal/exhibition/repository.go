package exhibition

import "context"

// Repository は Exhibition の保存・取得方法を表す境界。
//
// # Exhibition側は、保存先が
//
// - PostgreSQL
// - MySQL
// - SQLite
// - メモリ
// - 外部API
//
// のどれなのかを知らない。
//
// 「Exhibitionを保存できる」という能力だけを定義する。
type Repository interface {
	Save(
		ctx context.Context,
		exhibition Exhibition,
	) error
}
