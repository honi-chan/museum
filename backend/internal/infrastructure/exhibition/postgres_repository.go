package exhibition

import (
	"context"
	"errors"

	"museum/internal/domain"
	"museum/internal/infrastructure/database/sqlcgen"
	"museum/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewPostgresRepository(
	pool *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

// Save はExhibitionをPostgreSQLへ保存する。
func (r *PostgresRepository) Save(
	ctx context.Context,
	exhibition domain.Exhibition,
) error {
	return r.queries.CreateExhibition(
		ctx,
		sqlcgen.CreateExhibitionParams{
			ID:          exhibition.ID,
			MuseumID:    exhibition.MuseumID,
			Title:       exhibition.Title,
			Description: exhibition.Description,
			CreatedAt:   exhibition.CreatedAt,
		},
	)
}

// FindByID はIDを指定してExhibitionを取得する。
func (r *PostgresRepository) FindByID(
	ctx context.Context,
	id string,
) (domain.Exhibition, error) {
	row, err :=
		r.queries.GetExhibitionByID(
			ctx,
			id,
		)

	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return domain.Exhibition{},
				repository.ErrNotFound
		}

		return domain.Exhibition{}, err
	}

	return domain.Exhibition{
		ID:           row.ID,
		MuseumID:     row.MuseumID,
		Title:        row.Title,
		Description:  row.Description,
		DisplayOrder: row.DisplayOrder,
		CreatedAt:    row.CreatedAt,
	}, nil
}

// ListByMuseumID はMuseumに所属するExhibitionを
// display_order順で取得する。
func (r *PostgresRepository) ListByMuseumID(
	ctx context.Context,
	museumID string,
) ([]domain.Exhibition, error) {
	rows, err :=
		r.queries.ListExhibitionsByMuseumID(
			ctx,
			museumID,
		)

	if err != nil {
		return nil, err
	}

	exhibitions :=
		make(
			[]domain.Exhibition,
			0,
			len(rows),
		)

	for _, row := range rows {
		exhibitions = append(
			exhibitions,
			domain.Exhibition{
				ID:           row.ID,
				MuseumID:     row.MuseumID,
				Title:        row.Title,
				Description:  row.Description,
				DisplayOrder: row.DisplayOrder,
				CreatedAt:    row.CreatedAt,
			},
		)
	}

	return exhibitions, nil
}

var _ repository.ExhibitionRepository = (*PostgresRepository)(nil)
