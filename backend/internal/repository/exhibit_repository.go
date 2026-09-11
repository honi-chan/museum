package repository

import (
	"context"
	"museum/internal/domain"
)

type ExhibitRepository interface {
	// Save returns the persisted order, which is assigned atomically by storage.
	Save(context.Context, domain.Exhibit) (domain.Exhibit, error)
	ListByExhibitionID(context.Context, string) ([]domain.Exhibit, error)
}
