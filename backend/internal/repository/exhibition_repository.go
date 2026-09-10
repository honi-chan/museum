package repository

import (
	"context"

	"museum/internal/domain"
)

// ExhibitionRepository は
// Exhibitionの永続化に必要な操作を定義する。
//
// PostgreSQL / Memoryなどの
// 具体的な保存方式はここでは知らない。
type ExhibitionRepository interface {
	// Save はExhibitionを保存する。
	Save(
		ctx context.Context,
		exhibition domain.Exhibition,
	) error

	// FindByID はIDを指定して
	// Exhibitionを1件取得する。
	//
	// DBがPostgreSQLなのかMemoryなのかは
	// UseCaseから意識させない。
	FindByID(
		ctx context.Context,
		id string,
	) (domain.Exhibition, error)
}
