package exhibition

import (
	"context"
	"errors"
	"testing"
)

// errRepositoryUnavailable は、
// Repository保存失敗を再現するためのテスト用エラー。
//
// 実際のDB接続エラーなどを使う必要はなく、
// この固定エラーを使って
// 「保存失敗が呼び出し元まで返るか」を確認する。
var errRepositoryUnavailable = errors.New("repository unavailable")

// memoryRepository はテスト専用Repository。
//
// PostgreSQLを起動せずに、
// Repository.Save()が呼ばれたことを確認できる。
type memoryRepository struct {
	saved Exhibition
}

// Save は渡されたExhibitionを
// メモリ上に保持するだけ。
func (r *memoryRepository) Save(
	ctx context.Context,
	exhibition Exhibition,
) error {
	r.saved = exhibition

	return nil
}

func TestCreateAndSave(t *testing.T) {
	ctx := context.Background()

	input := CreateInput{
		MuseumID:    "museum-001",
		Title:       "THINGS I MADE",
		Description: "Things I created.",
	}

	// テストではIDを固定する。
	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	// DBの代わりにメモリRepositoryを使用する。
	repository := &memoryRepository{}

	created, err := CreateAndSave(
		ctx,
		input,
		idGenerator,
		repository,
	)

	if err != nil {
		t.Fatalf(
			"CreateAndSave() returned error: %v",
			err,
		)
	}

	// 作成されたExhibitionのIDを確認する。
	if created.ID != "exhibition-001" {
		t.Errorf(
			"ID = %q, want %q",
			created.ID,
			"exhibition-001",
		)
	}

	// Repositoryに保存された内容を確認する。
	if repository.saved.ID != created.ID {
		t.Errorf(
			"saved ID = %q, want %q",
			repository.saved.ID,
			created.ID,
		)
	}

	if repository.saved.Title != "THINGS I MADE" {
		t.Errorf(
			"saved Title = %q, want %q",
			repository.saved.Title,
			"THINGS I MADE",
		)
	}
}

// failingRepository は、
// 必ず保存に失敗するテスト用Repository。
type failingRepository struct{}

func (f failingRepository) Save(
	ctx context.Context,
	exhibition Exhibition,
) error {
	return errRepositoryUnavailable
}

func TestCreateAndSave_ReturnsErrorWhenSaveFails(
	t *testing.T,
) {
	ctx := context.Background()

	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "THINGS I MADE",
	}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	repository := failingRepository{}

	_, err := CreateAndSave(
		ctx,
		input,
		idGenerator,
		repository,
	)

	// Repository保存失敗が、
	// 呼び出し元まで返されることを確認する。
	if err == nil {
		t.Fatal(
			"CreateAndSave() expected error, got nil",
		)
	}

	// errors.Isを使って、
	// Repositoryが返したエラーがそのまま伝播しているか確認する。
	if !errors.Is(
		err,
		errRepositoryUnavailable,
	) {
		t.Errorf(
			"error = %v, want %v",
			err,
			errRepositoryUnavailable,
		)
	}
}
