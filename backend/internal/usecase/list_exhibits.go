package usecase

import (
	"context"
	"museum/internal/domain"
	"museum/internal/repository"
)

type ListExhibitsUseCase struct {
	exhibits    repository.ExhibitRepository
	exhibitions repository.ExhibitionRepository
}

func NewListExhibitsUseCase(exhibits repository.ExhibitRepository, exhibitions repository.ExhibitionRepository) *ListExhibitsUseCase {
	return &ListExhibitsUseCase{exhibits: exhibits, exhibitions: exhibitions}
}

func (u *ListExhibitsUseCase) Execute(ctx context.Context, exhibitionID string) ([]domain.Exhibit, error) {
	// An absent room is different from an existing room with no exhibits.
	if _, err := u.exhibitions.FindByID(ctx, exhibitionID); err != nil {
		return nil, err
	}
	return u.exhibits.ListByExhibitionID(ctx, exhibitionID)
}
