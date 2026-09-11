package repository

import (
	"context"

	"museum/internal/domain"
)

// ExhibitionRepository は
// Exhibitionの永続化に必要な操作を定義する。
//
// PostgreSQLなど具体的な保存方式は知らない。
type ExhibitionRepository interface {
	// Save はExhibitionを保存する。
	Save(
		ctx context.Context,
		exhibition domain.Exhibition,
	) error

	// FindByID はIDからExhibitionを取得する。
	FindByID(
		ctx context.Context,
		id string,
	) (domain.Exhibition, error)

	// ListByMuseumID は指定Museumに所属する
	// Exhibitionを表示順で取得する。
	ListByMuseumID(
		ctx context.Context,
		museumID string,
	) ([]domain.Exhibition, error)
}
