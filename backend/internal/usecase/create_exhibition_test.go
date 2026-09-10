package usecase

import (
	"context"
	"errors"
	"testing"

	"museum/internal/domain"
)

// --------------------------------------------------
// Test Double: IDGenerator
// --------------------------------------------------

// fixedIDGenerator はテスト専用のIDGenerator。
//
// UUIDのように毎回値が変わるIDを使うと
// テスト結果も毎回変わってしまうため、
// 固定値を返す実装を使う。
type fixedIDGenerator struct {
	id string
}

// Generate は常に固定されたIDを返す。
func (g fixedIDGenerator) Generate() string {
	return g.id
}

// --------------------------------------------------
// Test Double: Repository
// --------------------------------------------------

// fakeExhibitionRepository は
// テスト専用のRepository実装。
//
// PostgreSQLなどの本物のDBは使わず、
// Save()されたExhibitionをメモリ上に保持する。
type fakeExhibitionRepository struct {
	saved domain.Exhibition
	err   error
}

// Save は保存されたExhibitionを記録する。
//
// errが設定されている場合は、
// DB障害などを想定してそのerrorを返す。
func (r *fakeExhibitionRepository) Save(
	ctx context.Context,
	exhibition domain.Exhibition,
) error {
	if r.err != nil {
		return r.err
	}

	r.saved = exhibition

	return nil
}

// --------------------------------------------------
// 正常系
// --------------------------------------------------

func TestCreateExhibitionUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	// RepositoryはDBを使わずFakeを使用する。
	repository := &fakeExhibitionRepository{}

	// ID生成結果も固定する。
	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	// テスト対象のUseCaseを生成する。
	useCase := NewCreateExhibitionUseCase(
		repository,
		idGenerator,
	)

	input := CreateExhibitionInput{
		MuseumID:    "museum-001",
		Title:       "THINGS I MADE",
		Description: "Things I created.",
	}

	created, err := useCase.Execute(
		ctx,
		input,
	)

	if err != nil {
		t.Fatalf(
			"Execute() returned error: %v",
			err,
		)
	}

	// IDGeneratorが生成したIDが
	// Exhibitionに設定されることを確認する。
	if created.ID != "exhibition-001" {
		t.Errorf(
			"ID = %q, want %q",
			created.ID,
			"exhibition-001",
		)
	}

	// MuseumIDが正しく設定されることを確認する。
	if created.MuseumID != "museum-001" {
		t.Errorf(
			"MuseumID = %q, want %q",
			created.MuseumID,
			"museum-001",
		)
	}

	// Titleが正しく設定されることを確認する。
	if created.Title != "THINGS I MADE" {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			"THINGS I MADE",
		)
	}

	// Descriptionが正しく設定されることを確認する。
	if created.Description != "Things I created." {
		t.Errorf(
			"Description = %q, want %q",
			created.Description,
			"Things I created.",
		)
	}

	// CreatedAtが設定されていることを確認する。
	if created.CreatedAt.IsZero() {
		t.Error(
			"CreatedAt should not be zero",
		)
	}

	// Repositoryにも同じExhibitionが
	// 保存されていることを確認する。
	if repository.saved.ID != created.ID {
		t.Errorf(
			"saved ID = %q, want %q",
			repository.saved.ID,
			created.ID,
		)
	}
}

// --------------------------------------------------
// タイトル正規化
// --------------------------------------------------

func TestCreateExhibitionUseCase_Execute_TrimsTitle(
	t *testing.T,
) {
	ctx := context.Background()

	repository := &fakeExhibitionRepository{}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	useCase := NewCreateExhibitionUseCase(
		repository,
		idGenerator,
	)

	input := CreateExhibitionInput{
		MuseumID: "museum-001",

		// ユーザー入力に余計な空白が含まれているケース。
		Title: "   THINGS I MADE   ",
	}

	created, err := useCase.Execute(
		ctx,
		input,
	)

	if err != nil {
		t.Fatalf(
			"Execute() returned error: %v",
			err,
		)
	}

	// MUSEUM内部では前後の空白を除去した
	// 正規化済みタイトルを保持する。
	if created.Title != "THINGS I MADE" {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			"THINGS I MADE",
		)
	}
}

// --------------------------------------------------
// Validation
// --------------------------------------------------

func TestCreateExhibitionUseCase_Execute_RequiresTitle(
	t *testing.T,
) {
	ctx := context.Background()

	repository := &fakeExhibitionRepository{}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	useCase := NewCreateExhibitionUseCase(
		repository,
		idGenerator,
	)

	input := CreateExhibitionInput{
		MuseumID: "museum-001",
		Title:    "",
	}

	_, err := useCase.Execute(
		ctx,
		input,
	)

	// MUSEUMのルールとして、
	// タイトルのない展示室は作成できない。
	if err == nil {
		t.Fatal(
			"Execute() expected error, got nil",
		)
	}
}

func TestCreateExhibitionUseCase_Execute_RejectsWhitespaceOnlyTitle(
	t *testing.T,
) {
	ctx := context.Background()

	repository := &fakeExhibitionRepository{}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	useCase := NewCreateExhibitionUseCase(
		repository,
		idGenerator,
	)

	input := CreateExhibitionInput{
		MuseumID: "museum-001",

		// Trimすると空文字になる。
		Title: "     ",
	}

	_, err := useCase.Execute(
		ctx,
		input,
	)

	if err == nil {
		t.Fatal(
			"Execute() expected error, got nil",
		)
	}
}

// --------------------------------------------------
// Repository Error
// --------------------------------------------------

func TestCreateExhibitionUseCase_Execute_ReturnsRepositoryError(
	t *testing.T,
) {
	ctx := context.Background()

	// DB障害などを想定したエラー。
	expectedErr := errors.New(
		"repository unavailable",
	)

	repository := &fakeExhibitionRepository{
		err: expectedErr,
	}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	useCase := NewCreateExhibitionUseCase(
		repository,
		idGenerator,
	)

	input := CreateExhibitionInput{
		MuseumID: "museum-001",
		Title:    "THINGS I MADE",
	}

	_, err := useCase.Execute(
		ctx,
		input,
	)

	if err == nil {
		t.Fatal(
			"Execute() expected error, got nil",
		)
	}

	// Repositoryから返されたerrorが
	// 失われずそのまま返ることを確認する。
	if !errors.Is(
		err,
		expectedErr,
	) {
		t.Errorf(
			"error = %v, want %v",
			err,
			expectedErr,
		)
	}
}
