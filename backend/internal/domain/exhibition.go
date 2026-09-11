package domain

import "time"

// Exhibition はMUSEUMにおける「展示室」を表す。
//
// HTTP、DB、JSONなどの技術的な事情を持たず、
// MUSEUMそのものの概念だけを表現する。
type Exhibition struct {
	ID          string
	MuseumID    string
	Title       string
	Description string

	// DisplayOrder はMuseum内での表示順。
	//
	// MUSEUMでは作成日時ではなく、
	// ユーザー自身が決めた展示順を重要視する。
	DisplayOrder int32

	CreatedAt time.Time
}
