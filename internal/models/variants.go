package models

import "encoding/json"

// Variant represents the `variants` table.
// Many columns are nullable, hence pointer types are used.
type Variant struct {
    ID          int       `json:"id"`
    ProductID   *int      `json:"product_id"`
    VariantName *string   `json:"variant_name"`
    Pics        *[]string `json:"pics"`
    VariantType *string   `json:"variant_type"`
    IsWishlisted bool `json:"is_wishlisted"`
}

func (Variant) TableName() string { return "variants" }

// ScanPics unmarshals the jsonb pics column if the repository scans it raw.
func (v *Variant) ScanPics(raw []byte) error {
    if raw == nil {
        return nil
    }
    var pics []string
    if err := json.Unmarshal(raw, &pics); err != nil {
        return err
    }
    v.Pics = &pics
    return nil
} 