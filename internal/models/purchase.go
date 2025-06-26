package models

import (
    "time"

    "github.com/google/uuid"
)

// Purchase represents purchase_referral table.
type Purchase struct {
    ID         uuid.UUID `json:"id"`
    UserID     int       `json:"user_id"`
    DeliveryID uuid.UUID `json:"delivery_id"`
    CartID     uuid.UUID `json:"cart_id"`
    CreatedAt  time.Time `json:"created_at"`
}

func (Purchase) TableName() string { return "purchase_referral" } 