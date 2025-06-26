package models

import (
    "time"

    "github.com/google/uuid"
)

// DeliveryDetail represents a row in delivery_details_referral table.
// Nullable fields mapped as pointers for JSON omitempty.
type DeliveryDetail struct {
    ID                      uuid.UUID              `json:"id"`
    Name                    string                 `json:"name"`
    UserID                  int                    `json:"user_id"`
    DistrictID              int                    `json:"district_id"`
    ThanaID                 int                    `json:"thana_id"`
    FullAddress             string                 `json:"full_address"`
    PhoneNumber             string                 `json:"phone_number"`
    SecondaryPhoneNumber    *string                `json:"secondary_phone_number,omitempty"`
    Email                   *string                `json:"email,omitempty"`
    Data                    map[string]interface{} `json:"data,omitempty"`
    CreatedAt               time.Time              `json:"created_at"`
    UpdatedAt               time.Time              `json:"updated_at"`
    PrimaryHasWhatsApp      bool                   `json:"primary_number_has_whatsapp"`
    SecondaryHasWhatsApp    bool                   `json:"secondary_number_has_whatsapp"`
}

func (DeliveryDetail) TableName() string { return "delivery_details_referral" } 