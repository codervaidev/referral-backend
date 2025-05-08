package models

// ReferralUser represents the referral_user table in the database
type ReferralUser struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	RefererCode string `json:"referer_code"`
	Points      int    `json:"points"`
	User        User   `json:"user"`
}

// TableName returns the table name for ReferralUser
func (ReferralUser) TableName() string {
	return "referral_user"
}
