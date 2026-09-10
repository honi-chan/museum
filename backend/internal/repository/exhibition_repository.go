package repository

import (
	"context"

	"museum/internal/domain"
)

// ExhibitionRepository は
// Exhibitionの永続化に必要な操作を定義する。
//
// 実際の保存先が
//
// - PostgreSQL
// - MySQL
// - Memory
// - 外部API
//
// のどれなのかは、このinterfaceでは知らない。
type ExhibitionRepository interface {
	Save(
		ctx context.Context,
		exhibition domain.Exhibition,
	) error
}
