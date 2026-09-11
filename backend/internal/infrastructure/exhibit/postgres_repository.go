package exhibit

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"museum/internal/domain"
	"museum/internal/infrastructure/database/sqlcgen"
	"museum/internal/repository"
)

type PostgresRepository struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool, queries: sqlcgen.New(pool)}
}

func (r *PostgresRepository) Save(ctx context.Context, exhibit domain.Exhibit) (domain.Exhibit, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return domain.Exhibit{}, err
	}
	defer tx.Rollback(ctx)
	queries := r.queries.WithTx(tx)
	// Lock the parent even for an empty room. The separate INSERT sees the latest
	// committed order after waiting, so concurrent writers append sequentially.
	if _, err = queries.LockExhibitionForExhibit(ctx, exhibit.ExhibitionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Exhibit{}, repository.ErrNotFound
		}
		return domain.Exhibit{}, err
	}
	row, err := queries.CreateExhibit(ctx, sqlcgen.CreateExhibitParams{
		ID: exhibit.ID, ExhibitionID: exhibit.ExhibitionID, Title: exhibit.Title,
		Caption: exhibit.Caption, ImageUrl: exhibit.ImageURL, CreatedAt: exhibit.CreatedAt,
	})
	if err != nil {
		return domain.Exhibit{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Exhibit{}, err
	}
	return toDomain(row), nil
}

func (r *PostgresRepository) ListByExhibitionID(ctx context.Context, id string) ([]domain.Exhibit, error) {
	rows, err := r.queries.ListExhibitsByExhibitionID(ctx, id)
	if err != nil {
		return nil, err
	}
	exhibits := make([]domain.Exhibit, 0, len(rows))
	for _, row := range rows {
		exhibits = append(exhibits, toDomain(row))
	}
	return exhibits, nil
}

func toDomain(row sqlcgen.Exhibit) domain.Exhibit {
	return domain.Exhibit{ID: row.ID, ExhibitionID: row.ExhibitionID, Title: row.Title,
		Caption: row.Caption, ImageURL: row.ImageUrl, DisplayOrder: row.DisplayOrder, CreatedAt: row.CreatedAt}
}

var _ repository.ExhibitRepository = (*PostgresRepository)(nil)
