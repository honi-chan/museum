package usecase

import (
	"context"

	"museum/internal/domain"
	"museum/internal/repository"
)

// GetExhibitionUseCase は
// 「展示室を1件取得する」ユースケース。
type GetExhibitionUseCase struct {
	repository repository.ExhibitionRepository
}

// NewGetExhibitionUseCase は
// Repositoryを注入してUseCaseを生成する。
func NewGetExhibitionUseCase(
	repository repository.ExhibitionRepository,
) *GetExhibitionUseCase {
	return &GetExhibitionUseCase{
		repository: repository,
	}
}

// Execute はIDを指定して
// Exhibitionを取得する。
//
// HTTPやPostgreSQLについては知らない。
func (u *GetExhibitionUseCase) Execute(
	ctx context.Context,
	id string,
) (domain.Exhibition, error) {
	return u.repository.FindByID(
		ctx,
		id,
	)
}
