package models

import "time"

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

// TableName returns the table name for User
func (User) TableName() string {
	return "User"
} 