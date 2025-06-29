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

// GetByID fetches delivery details by ID
func (r *DeliveryDetailsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.DeliveryDetail, error) {
    const q = `SELECT id, name, user_id, district_id, thana_id, full_address, phone_number, secondary_phone_number, email, data, created_at, updated_at, primary_number_has_whatsapp, secondary_number_has_whatsapp 
               FROM delivery_details_referral WHERE id = $1`
    
    var delivery models.DeliveryDetail
    var dataJSON []byte
    
    err := r.DB.QueryRow(ctx, q, id).Scan(
        &delivery.ID, &delivery.Name, &delivery.UserID, &delivery.DistrictID, &delivery.ThanaID,
        &delivery.FullAddress, &delivery.PhoneNumber, &delivery.SecondaryPhoneNumber, &delivery.Email,
        &dataJSON, &delivery.CreatedAt, &delivery.UpdatedAt, &delivery.PrimaryHasWhatsApp, &delivery.SecondaryHasWhatsApp,
    )
    if err != nil {
        return nil, err
    }
    
    if len(dataJSON) > 0 {
        if err := json.Unmarshal(dataJSON, &delivery.Data); err != nil {
            return nil, err
        }
    }
    
    return &delivery, nil
} 