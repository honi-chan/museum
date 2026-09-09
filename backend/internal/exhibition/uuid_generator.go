package exhibition

import "github.com/google/uuid"

// UUIDGenerator は本番用のIDGenerator。
//
// Exhibition側からは、
// UUIDを使用していることは見えない。
type UUIDGenerator struct{}

// Generate は新しいUUIDを生成する。
func (UUIDGenerator) Generate() string {
	return uuid.NewString()
}
