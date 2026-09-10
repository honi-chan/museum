package exhibition

import (
	"context"
	"errors"

	"museum/internal/domain"
	"museum/internal/infrastructure/exhibition/sqlcgen"
	"museum/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository は
// ExhibitionRepositoryのPostgreSQL実装。
//
// SQLそのものはsqlcへ任せ、
// ここでは
//
// DB Model
// ↓
// Domain
//
// の変換を担当する。
type PostgresRepository struct {
	queries *sqlcgen.Queries
}

// NewPostgresRepository は
// PostgreSQL Repositoryを生成する。
func NewPostgresRepository(
	pool *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

// Save はExhibitionを保存する。
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

// FindByID はExhibitionをIDで取得する。
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
		// PostgreSQL / pgx固有の
		// ErrNoRowsをRepository共通エラーへ変換する。
		//
		// UseCaseにpgx依存を漏らさない。
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return domain.Exhibition{},
				repository.ErrNotFound
		}

		return domain.Exhibition{}, err
	}

	// sqlc生成Modelを
	// Domain Modelへ変換する。
	return domain.Exhibition{
		ID:          row.ID,
		MuseumID:    row.MuseumID,
		Title:       row.Title,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
	}, nil
}

// Compile-time check.
//
// PostgresRepositoryが
// ExhibitionRepositoryを実装していることを
// コンパイル時に保証する。
var _ repository.ExhibitionRepository = (*PostgresRepository)(nil)
