package exhibition

import "github.com/google/uuid"

// UUIDGenerator は
// UUIDを使ったID生成の具体実装。
//
// UseCaseはこの具体実装を知らない。
type UUIDGenerator struct{}

// Generate はUUIDを生成する。
func (UUIDGenerator) Generate() string {
	return uuid.NewString()
}
