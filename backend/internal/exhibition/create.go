package exhibition

import (
	"errors"
	"strings"
	"time"
)

// CreateInput は、展示室を作成するために必要な情報。
//
// HTTP Requestではない。
// Web / CLI / Mobileなど、
// 呼び出し元に依存しない入力モデル。
type CreateInput struct {
	MuseumID    string
	Title       string
	Description string
}

// Create は新しい展示室を生成する。
//
// Exhibitionを作るためのルールだけを担当する。
//
// IDを「どう生成するか」はCreateの責務ではないため、
// IDGeneratorとして外から受け取る。
func Create(
	input CreateInput,
	idGenerator IDGenerator,
) (Exhibition, error) {
	// タイトル前後の余計な空白を除去する。
	title := strings.TrimSpace(input.Title)

	// MUSEUMのルール:
	// タイトルのない展示室は作成できない。
	if title == "" {
		return Exhibition{}, errors.New("exhibition title is required")
	}

	// ID生成方法そのものは知らない。
	//
	// Createが知っているのは、
	// Generate()を呼べばIDが取得できることだけ。
	id := idGenerator.Generate()

	return Exhibition{
		ID:          id,
		MuseumID:    input.MuseumID,
		Title:       title,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}, nil
}
