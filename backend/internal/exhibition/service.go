package exhibition

import "context"

// Service は Exhibition に関する操作をまとめる。
//
// Handler側が Repository や IDGenerator を
// 毎回意識しなくて済むようにする。
//
// 将来、依存関係が増えても
// Serviceの生成部分だけ直せばよい。
type Service struct {
	repository  Repository
	idGenerator IDGenerator
}

// NewService は Exhibition Service を生成する。
//
// 依存関係をここで一度だけ受け取る。
func NewService(
	repository Repository,
	idGenerator IDGenerator,
) *Service {
	return &Service{
		repository:  repository,
		idGenerator: idGenerator,
	}
}

// Create は展示室作成ユースケースを実行する。
//
// Handlerは入力だけ渡せばよく、
// RepositoryやID生成方法を知る必要がない。
func (s *Service) Create(
	ctx context.Context,
	input CreateInput,
) (Exhibition, error) {
	return CreateAndSave(
		ctx,
		input,
		s.idGenerator,
		s.repository,
	)
}
