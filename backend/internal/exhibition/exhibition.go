package exhibition

import "time"

// Exhibition は、MUSEUM内に存在する「展示室」を表す。
//
// 重要:
// Exhibition はSNSのカテゴリやタグではない。
// ユーザーが意図を持って作品をまとめる「空間」である。
//
// 例:
// - THINGS I MADE
// - SUMMER IN KYOTO
// - THINGS I LOVED
type Exhibition struct {
	ID          string
	MuseumID    string
	Title       string
	Description string
	CreatedAt   time.Time
}
