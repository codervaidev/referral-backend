package models

import (
	"time"

	"github.com/google/uuid"
)

type GemHistory struct {
	ID         int       `json:"id"`
	Amount     int       `json:"amount"`
	UserID     int       `json:"user_id"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`
	PurchaseID *uuid.UUID `json:"purchase_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type GemHistoryResponse struct {
	ID         int       `json:"id"`
	Amount     int       `json:"amount"`
	UserID     int       `json:"user_id"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`
	PurchaseID *uuid.UUID `json:"purchase_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Phone      string    `json:"phone"`
	Name       string    `json:"name"`
	ImageUrl   *string   `json:"image_url"`
	// Purchase details for purchase type
	PurchaseDetails *PurchaseHistoryDetails `json:"purchase_details,omitempty"`
}

// PurchaseHistoryDetails contains detailed purchase information
type PurchaseHistoryDetails struct {
	PurchaseID       uuid.UUID              `json:"purchase_id"`
	DeliveryLocation string                 `json:"delivery_location"`
	DeliveryFee      int                   `json:"delivery_fee"`
	TotalAmount      int                   `json:"total_amount"`
	Items            []PurchaseHistoryItem  `json:"items"`
}

type PurchaseHistoryItem struct {
	ProductID   int     `json:"product_id"`
	ProductName *string `json:"product_name"`
	VariantID   *int    `json:"variant_id,omitempty"`
	VariantName *string `json:"variant_name,omitempty"`
	Quantity    int     `json:"quantity"`
	PriceAtAdd  float64 `json:"price_at_add"`
}

func (GemHistory) TableName() string {
	return "gems_history"
}
