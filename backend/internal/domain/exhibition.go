package domain

import "time"

// Exhibition はMUSEUMにおける「展示室」を表す。
//
// HTTP、DB、JSON、UUIDなどの
// 技術的な事情を一切持たない。
//
// ここにはMUSEUMそのものの概念だけを置く。
type Exhibition struct {
	ID          string
	MuseumID    string
	Title       string
	Description string
	CreatedAt   time.Time
}
