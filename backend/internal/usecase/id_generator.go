package usecase

// IDGenerator はID生成機能の境界。
//
// UUIDを使うのかULIDを使うのかを
// UseCaseは知らなくてよい。
type IDGenerator interface {
	Generate() string
}
