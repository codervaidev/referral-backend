package models

type GemHistory struct {
	ID      int    `json:"id"`
	Amount  int    `json:"amount"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (GemHistory) TableName() string {
	return "gems_history"
}
