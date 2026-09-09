package exhibition

import (
	"errors"
	"strings"
	"time"
)

// CreateInput は、展示室を作成するために必要な情報。
//
// HTTPのRequestではない。
// CLIやモバイルアプリなど、どこから呼ばれても使える入力モデル。
type CreateInput struct {
	MuseumID    string
	Title       string
	Description string
}

// Create は新しい展示室を生成する。
//
// この関数には「展示室を作る」という
// MUSEUM固有のルールを書く。
func Create(input CreateInput) (Exhibition, error) {
	// 展示室には必ずタイトルが必要。
	title := strings.TrimSpace(input.Title)

	if title == "" {
		return Exhibition{}, errors.New("exhibition title is required")
	}

	// 現時点ではID生成方法は仮。
	//
	// 後からUUID生成処理などが必要になったとき、
	// この部分の責務を分離する。
	id := "temporary-id"

	return Exhibition{
		ID:          id,
		MuseumID:    input.MuseumID,
		Title:       title,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}, nil
}
