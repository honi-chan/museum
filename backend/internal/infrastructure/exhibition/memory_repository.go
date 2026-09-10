package exhibition

import (
	"context"
	"sync"

	"museum/internal/domain"
)

// MemoryRepository は
// Exhibitionをメモリに保存する具体実装。
//
// PostgreSQL実装前の開発・テスト用途として利用する。
type MemoryRepository struct {
	mu sync.Mutex

	exhibitions []domain.Exhibition
}

// NewMemoryRepository は
// MemoryRepositoryを生成する。
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		exhibitions: make(
			[]domain.Exhibition,
			0,
		),
	}
}

// Save はExhibitionをメモリへ保存する。
func (r *MemoryRepository) Save(
	ctx context.Context,
	exhibition domain.Exhibition,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.exhibitions = append(
		r.exhibitions,
		exhibition,
	)

	return nil
}
