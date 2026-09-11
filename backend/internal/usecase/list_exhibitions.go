package usecase

import (
	"context"

	"museum/internal/domain"
	"museum/internal/repository"
)

// ListExhibitionsUseCase は
// Museum内の展示室一覧を取得するユースケース。
type ListExhibitionsUseCase struct {
	repository repository.ExhibitionRepository
}

// NewListExhibitionsUseCase は
// Repositoryを注入してUseCaseを生成する。
func NewListExhibitionsUseCase(
	repository repository.ExhibitionRepository,
) *ListExhibitionsUseCase {
	return &ListExhibitionsUseCase{
		repository: repository,
	}
}

// Execute はMuseum IDを指定して
// Exhibition一覧を取得する。
func (u *ListExhibitionsUseCase) Execute(
	ctx context.Context,
	museumID string,
) ([]domain.Exhibition, error) {
	return u.repository.ListByMuseumID(
		ctx,
		museumID,
	)
}
