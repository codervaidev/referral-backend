package models

import "github.com/google/uuid"

type Gem struct {
    ID          uuid.UUID `json:"id"`
    Name        *string   `json:"name"`
    Image       *string   `json:"image"`
    GemsCount   *int      `json:"gems_count"`
	IsActive    bool      `json:"is_active"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Category    *string   `json:"category"`
    Link        *string   `json:"link"`
    Variant     *string   `json:"variant"`
    Picture1    *string   `json:"picture_1"`
    Picture2    *string   `json:"picture_2"`
    Picture3    *string   `json:"picture_3"`
    Picture4    *string   `json:"picture_4"`
}

func (Gem) TableName() string {
	return "gems_store"
}
