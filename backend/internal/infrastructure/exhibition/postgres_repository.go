package exhibition

import (
	"context"

	"museum/internal/domain"
	"museum/internal/infrastructure/exhibition/sqlcgen"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository は
// ExhibitionRepositoryのPostgreSQL実装。
//
// UseCaseからPostgreSQL固有のコードを
// 完全に隔離する。
type PostgresRepository struct {
	queries *sqlcgen.Queries
}

// NewPostgresRepository は
// PostgreSQL Repositoryを生成する。
func NewPostgresRepository(
	pool *pgxpool.Pool,
) *PostgresRepository {
	return &PostgresRepository{
		// sqlcが生成したQueriesに
		// pgx connection poolを渡す。
		queries: sqlcgen.New(pool),
	}
}

// Save はExhibitionをPostgreSQLへ保存する。
//
// INSERT文そのものはここに書かない。
// SQLはdb/query/exhibitions.sqlで管理し、
// Goコードはsqlcが生成する。
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
