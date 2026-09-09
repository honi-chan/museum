package exhibition

import "testing"

// 展示室を正常に作成できることを確認する。
func TestCreate(t *testing.T) {
	input := CreateInput{
		MuseumID:    "museum-001",
		Title:       "THINGS I MADE",
		Description: "Things I created.",
	}

	created, err := Create(input)

	// エラーが返ってきたらテスト失敗。
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	// MuseumIDが入力値と一致すること。
	if created.MuseumID != input.MuseumID {
		t.Errorf(
			"MuseumID = %q, want %q",
			created.MuseumID,
			input.MuseumID,
		)
	}

	// Titleが入力値と一致すること。
	if created.Title != input.Title {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			input.Title,
		)
	}

	// Descriptionが入力値と一致すること。
	if created.Description != input.Description {
		t.Errorf(
			"Description = %q, want %q",
			created.Description,
			input.Description,
		)
	}

	// 作成された展示室にはIDが存在すること。
	if created.ID == "" {
		t.Error("ID should not be empty")
	}

	// 作成日時が設定されていること。
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

// タイトル前後の空白が除去されることを確認する。
func TestCreate_TrimsTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "  THINGS I MADE  ",
	}

	created, err := Create(input)

	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	// UIやAPIから余計な空白が渡されても、
	// Exhibition内部では正規化された状態にする。
	if created.Title != "THINGS I MADE" {
		t.Errorf(
			"Title = %q, want %q",
			created.Title,
			"THINGS I MADE",
		)
	}
}

// タイトルが空の場合は展示室を作れないことを確認する。
func TestCreate_RequiresTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "",
	}

	_, err := Create(input)

	// MUSEUMのルールとして、
	// タイトルのない展示室は存在できない。
	if err == nil {
		t.Fatal("Create() expected error, got nil")
	}
}

// 空白だけのタイトルも無効であることを確認する。
func TestCreate_RejectsWhitespaceOnlyTitle(t *testing.T) {
	input := CreateInput{
		MuseumID: "museum-001",
		Title:    "    ",
	}

	_, err := Create(input)

	if err == nil {
		t.Fatal("Create() expected error, got nil")
	}
}
