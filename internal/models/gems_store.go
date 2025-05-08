package models

import "github.com/google/uuid"

type Gem struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Image     string    `json:"image"`
    GemsCount int       `json:"gems_count"`
    IsActive  bool      `json:"is_active"`
}

func (Gem) TableName() string {
	return "gems_store"
}
