package repository

import (
	"context"
	"time"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDetailsRepo struct {
	DB *pgxpool.Pool
}

func NewUserDetailsRepo(db *pgxpool.Pool) *UserDetailsRepo {
	return &UserDetailsRepo{DB: db}
}

func (r *UserDetailsRepo) Create(ctx context.Context, details models.UserDetails) error {
	query := `
		INSERT INTO referral_details (
			id, user_id, name, district_id, thana_id, full_address, phone_number, secondary_phone_number, data, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query,
		details.ID,
		details.UserID,
		details.Name,
		details.DistrictID,
		details.ThanaID,
		details.FullAddress,
		details.PhoneNumber,
		details.SecondaryPhoneNumber,
		details.Data,
		now,
		now,
	)
	return err
}

func (r *UserDetailsRepo) GetByUserID(ctx context.Context, userID uint) (*models.UserDetails, error) {
	query := `
		SELECT id, user_id, name, district_id, thana_id, full_address, phone_number, secondary_phone_number, data, created_at, updated_at
		FROM referral_details
		WHERE user_id = $1
	`
	var details models.UserDetails
	err := r.DB.QueryRow(ctx, query, userID).Scan(
		&details.ID,
		&details.UserID,
		&details.Name,
		&details.DistrictID,
		&details.ThanaID,
		&details.FullAddress,
		&details.PhoneNumber,
		&details.SecondaryPhoneNumber,
		&details.Data,
		&details.CreatedAt,
		&details.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &details, nil
}

func (r *UserDetailsRepo) Update(ctx context.Context, details models.UserDetails) error {
	query := `
		UPDATE referral_details
		SET name = $1, district_id = $2, thana_id = $3, full_address = $4, phone_number = $5, secondary_phone_number = $6, data = $7, updated_at = $8
		WHERE user_id = $9
	`
	_, err := r.DB.Exec(ctx, query,
		details.Name,
		details.DistrictID,
		details.ThanaID,
		details.FullAddress,
		details.PhoneNumber,
		details.SecondaryPhoneNumber,
		details.Data,
		time.Now(),
		details.UserID,
	)
	return err
} 