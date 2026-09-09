package exhibition

import "testing"

// fixedIDGenerator はテスト専用のIDGenerator。
//
// Generate()を呼ぶと必ず同じIDを返す。
//
// 本物のUUIDなどを使わないことで、
// テスト結果を毎回同じにできる。
type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) Generate() string {
	return g.id
}

func TestCreate(t *testing.T) {
	input := CreateInput{
		MuseumID:    "museum-001",
		Title:       "THINGS I MADE",
		Description: "Things I created.",
	}

	// テストではID生成結果を固定する。
	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	created, err := Create(
		input,
		idGenerator,
	)

	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	// IDGeneratorが返したIDが、
	// Exhibitionに設定されること。
	if created.ID != "exhibition-001" {
		t.Errorf(
			"ID = %q, want %q",
			created.ID,
			"exhibition-001",
		)
	}

	if created.MuseumID != input.MuseumID {
		t.Errorf(
			"MuseumID = %q, want %q",
			created.MuseumID,
			input.MuseumID,
		)
	}

	if created.Title != input.Title {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			input.Title,
		)
	}

	if created.Description != input.Description {
		t.Errorf(
			"Description = %q, want %q",
			created.Description,
			input.Description,
		)
	}

	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestCreate_TrimsTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "  THINGS I MADE  ",
	}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	created, err := Create(
		input,
		idGenerator,
	)

	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if created.Title != "THINGS I MADE" {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			"THINGS I MADE",
		)
	}
}

func TestCreate_RequiresTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "",
	}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	_, err := Create(
		input,
		idGenerator,
	)

	if err == nil {
		t.Fatal("Create() expected error, got nil")
	}
}

func TestCreate_RejectsWhitespaceOnlyTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "    ",
	}

	idGenerator := fixedIDGenerator{
		id: "exhibition-001",
	}

	_, err := Create(
		input,
		idGenerator,
	)

	if err == nil {
		t.Fatal("Create() expected error, got nil")
	}
}
