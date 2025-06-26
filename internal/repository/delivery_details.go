package repository

import (
    "context"
    "encoding/json"

    "github.com/codervaidev/referral-backend/internal/models"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryDetailsRepo struct {
    DB *pgxpool.Pool
}

func NewDeliveryDetailsRepo(db *pgxpool.Pool) *DeliveryDetailsRepo { return &DeliveryDetailsRepo{DB: db} }

// Create inserts a delivery detail row and returns the generated id.
func (r *DeliveryDetailsRepo) Create(ctx context.Context, d models.DeliveryDetail) (uuid.UUID, error) {
    var id uuid.UUID
    var dataJSON []byte
    if d.Data != nil {
        dataJSON, _ = json.Marshal(d.Data)
    }
    const q = `INSERT INTO delivery_details_referral
        (name, user_id, district_id, thana_id, full_address, phone_number, secondary_phone_number, email, data, primary_number_has_whatsapp, secondary_number_has_whatsapp)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        RETURNING id`

    err := r.DB.QueryRow(ctx, q, d.Name, d.UserID, d.DistrictID, d.ThanaID, d.FullAddress, d.PhoneNumber, d.SecondaryPhoneNumber, d.Email, dataJSON, d.PrimaryHasWhatsApp, d.SecondaryHasWhatsApp).Scan(&id)
    return id, err
} 