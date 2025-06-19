package models

import "encoding/json"

// Product represents a row from the `products` table.
// Most columns are nullable in the DDL, hence pointer types are used so that
// we can distinguish between NULL and zero-values when scanning.
type Product struct {
    ID          int       `json:"id"`
    CategoryID  *int      `json:"category_id"`
    Link        *string   `json:"link"`
    Title       *string   `json:"title"`
    Description *string   `json:"description"`
    Price       *float64  `json:"price"`
    Stock       *int      `json:"stock"`
    Sold        *int      `json:"sold"`
    WishlistCount *int     `json:"wishlist_count"`
    Rating        *float64 `json:"rating"`
    RecommendedFor *[]int  `json:"recommended_for"`
    ImageURLs   *[]string `json:"image_urls"`
    Vendor      *string   `json:"vendor"`
    Variants    []Variant `json:"variants,omitempty"`
}

// TableName overrides the table name used by any ORM that relies on it.
func (Product) TableName() string {
    return "products"
}

// ScanImageURLs is a small helper that can be used by repositories when the
// `image_urls` column is returned as raw JSONB (i.e. []byte). It unmarshals
// the JSON array into the struct field.
func (p *Product) ScanImageURLs(raw []byte) error {
    if raw == nil {
        return nil
    }
    var urls []string
    if err := json.Unmarshal(raw, &urls); err != nil {
        return err
    }
    p.ImageURLs = &urls
    return nil
}
