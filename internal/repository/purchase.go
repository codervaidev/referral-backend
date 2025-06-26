package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type PurchaseRepo struct {
    DB *pgxpool.Pool
}

func NewPurchaseRepo(db *pgxpool.Pool) *PurchaseRepo { return &PurchaseRepo{DB: db} }

func (r *PurchaseRepo) Create(ctx context.Context, userID uint, deliveryID, cartID uuid.UUID) (uuid.UUID, error) {
    var id uuid.UUID
    const q = `INSERT INTO purchase_referral (user_id, delivery_id, cart_id) VALUES ($1,$2,$3) RETURNING id`
    err := r.DB.QueryRow(ctx, q, userID, deliveryID, cartID).Scan(&id)
    return id, err
} 