package models

import "time"

type GemHistory struct {
	ID        int    `json:"id"`
	Amount    int    `json:"amount"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type GemHistoryResponse struct {
	ID        int    `json:"id"`
	Amount    int    `json:"amount"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string `json:"phone"`
	Name      string `json:"name"`
	ImageUrl  string `json:"image_url"`
}

func (GemHistory) TableName() string {
	return "gems_history"
}
