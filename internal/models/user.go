package models

import (
	"time"
     "github.com/google/uuid"
)

// User represents the User table in the database
type User struct {
	ID        uint      `json:"id"`
	Phone     string    `json:"phone"`
	PIN       string    `json:"pin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ParentID  *uint     `json:"parent_id"`
	Type      string    `json:"type"`
}	

// UserDetails represents the user profile information
type UserDetails struct {
	ID                    uuid.UUID      `json:"id"`
	UserID                int            `json:"user_id"`
	Name                  string         `json:"name"`
	DistrictID            int            `json:"district_id"`
	ThanaID               int            `json:"thana_id"`
	FullAddress           string         `json:"full_address"`
	PhoneNumber           string         `json:"phone_number"`
	SecondaryPhoneNumber  *string        `json:"secondary_phone_number"`
	Data                  map[string]interface{} `json:"data"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "User"
}

// TableName returns the table name for UserDetails
func (UserDetails) TableName() string {
	return "referral_details"
} 