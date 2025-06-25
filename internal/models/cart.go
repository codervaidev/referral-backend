package models

import (
	"time"

	"github.com/google/uuid"
)

// Cart reflects the `carts` table. Items are loaded separately.
// All fields are non-pointer because the table declares NOT NULL (except timestamps which default to now()).
type Cart struct {
    ID        uuid.UUID     `json:"id"`
    UserID    int           `json:"user_id"`
    Status    string        `json:"status"`
    CreatedAt time.Time     `json:"created_at"`
    UpdatedAt time.Time     `json:"updated_at"`
    Items     []CartItem    `json:"items"`
}

// TableName satisfies any ORM helper that might use reflection.
func (Cart) TableName() string { return "carts" }

// CartItem represents a row in the `cart_items` table. Product and Variant
// details may be eagerly loaded by repository methods.
// The VariantID column is nullable, hence a pointer is used.
type CartItem struct {
    ID          uuid.UUID  `json:"id"`
    CartID      uuid.UUID  `json:"cart_id"`
    ProductID   int        `json:"product_id"`
    VariantID   *int       `json:"variant_id,omitempty"`
    Quantity    int        `json:"quantity"`
    PriceAtAdd  float64    `json:"price_at_add"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`

    // Optional pre-loaded relations
    Product     *Product   `json:"product,omitempty"`
    Variant     *Variant   `json:"variant,omitempty"`
}

func (CartItem) TableName() string { return "cart_items" } 