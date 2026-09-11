package exhibit

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"museum/internal/domain"
	"museum/internal/repository"
)

// An isolated schema exercises real migrations without touching development data.
func TestPostgresRepositoryIntegration(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is required for PostgreSQL integration")
	}
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close(ctx)
	schema := "exhibit_test_" + fmt.Sprintf("%d", time.Now().UnixNano())
	name := pgx.Identifier{schema}.Sanitize()
	if _, err = admin.Exec(ctx, "CREATE SCHEMA "+name); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec(ctx, "DROP SCHEMA "+name+" CASCADE")
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	migrations, err := filepath.Glob("../../../db/migrations/*.up.sql")
	if err != nil || len(migrations) != 3 {
		t.Fatalf("migrations: %v %v", migrations, err)
	}
	for _, file := range migrations {
		sql, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = pool.Exec(ctx, string(sql)); err != nil {
			t.Fatalf("apply %s: %v", file, err)
		}
	}
	if _, err = pool.Exec(ctx, "INSERT INTO exhibitions (id,museum_id,title,created_at) VALUES ('room','museum','Room',now()), ('other','museum','Other',now())"); err != nil {
		t.Fatal(err)
	}
	repo := NewPostgresRepository(pool)
	if _, err = repo.Save(ctx, domain.Exhibit{ID: uuid.NewString(), ExhibitionID: "missing", ImageURL: "https://example.com/a", CreatedAt: time.Now()}); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("missing parent: %v", err)
	}
	const count = 12
	results := make(chan domain.Exhibit, count)
	failures := make(chan error, count)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, err := repo.Save(ctx, domain.Exhibit{ID: uuid.NewString(), ExhibitionID: "room", ImageURL: "https://example.com/a", CreatedAt: time.Now()})
			if err != nil {
				failures <- err
			} else {
				results <- got
			}
		}()
	}
	wg.Wait()
	close(failures)
	close(results)
	for err := range failures {
		t.Fatal(err)
	}
	seen := map[int32]bool{}
	for got := range results {
		if seen[got.DisplayOrder] {
			t.Fatalf("duplicate order %d", got.DisplayOrder)
		}
		seen[got.DisplayOrder] = true
	}
	if len(seen) != count {
		t.Fatalf("saved %d", len(seen))
	}
	items, err := repo.ListByExhibitionID(ctx, "room")
	if err != nil || len(items) != count {
		t.Fatalf("list: %v %v", items, err)
	}
	for i, item := range items {
		if item.DisplayOrder != int32(i) {
			t.Fatalf("order: %+v", items)
		}
	}
	other, err := repo.Save(ctx, domain.Exhibit{ID: "other-work", ExhibitionID: "other", ImageURL: "https://example.com/a", CreatedAt: time.Now()})
	if err != nil || other.DisplayOrder != 0 {
		t.Fatalf("room isolation: %+v %v", other, err)
	}
	// Equal positions must still have a deterministic ID ordering.
	if _, err = pool.Exec(ctx, "UPDATE exhibits SET display_order=0 WHERE exhibition_id='room'"); err != nil {
		t.Fatal(err)
	}
	items, err = repo.ListByExhibitionID(ctx, "room")
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(items); i++ {
		if items[i-1].ID >= items[i].ID {
			t.Fatal("unstable ID order")
		}
	}
	if _, err = pool.Exec(ctx, "DELETE FROM exhibitions WHERE id='room'"); err != nil {
		t.Fatal(err)
	}
	items, err = repo.ListByExhibitionID(ctx, "room")
	if err != nil || len(items) != 0 {
		t.Fatalf("cascade: %v %v", items, err)
	}
	items, err = repo.ListByExhibitionID(ctx, "other")
	if err != nil || len(items) != 1 {
		t.Fatalf("unrelated room changed: %v %v", items, err)
	}
	down, err := os.ReadFile("../../../db/migrations/000003_create_exhibits.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, string(down)); err != nil {
		t.Fatal(err)
	}
	up, err := os.ReadFile(migrations[2])
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, string(up)); err != nil {
		t.Fatal(err)
	}
}
