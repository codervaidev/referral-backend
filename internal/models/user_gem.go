package models

// ReferralUser represents the referral_user table in the database
type ReferralUser struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	RefererCode string `json:"referer_code"`
	Points      int    `json:"points"`
	User        User   `json:"user"`
}

type RedeemResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Points  int    `json:"points"`
}

// RedeemReferralCodeRequest represents the request body for redeeming referral points
type RedeemReferralCodeRequest struct {
	Points int `json:"points"`
}

// TableName returns the table name for ReferralUser
func (ReferralUser) TableName() string {
	return "referral_user"
}

func (RedeemReferralCodeRequest) TableName() string {
	return "referral_user"
}
