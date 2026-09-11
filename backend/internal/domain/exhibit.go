package domain

import "time"

type Exhibit struct {
	ID           string
	ExhibitionID string
	Title        string
	Caption      string
	ImageURL     string
	DisplayOrder int32
	CreatedAt    time.Time
}
