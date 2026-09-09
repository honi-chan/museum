package exhibition

import (
	"context"
	"sync"
)

// MemoryRepository は、
// Exhibitionをメモリ上に保存するRepository。
//
// 現段階ではPostgreSQL実装前の仮実装。
//
// サーバーを再起動するとデータは消える。
type MemoryRepository struct {
	mu sync.Mutex

	exhibitions []Exhibition
}

// NewMemoryRepository は
// MemoryRepositoryを生成する。
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		exhibitions: make(
			[]Exhibition,
			0,
		),
	}
}

// Save はExhibitionをメモリへ保存する。
func (r *MemoryRepository) Save(
	ctx context.Context,
	exhibition Exhibition,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.exhibitions = append(
		r.exhibitions,
		exhibition,
	)

	return nil
}
