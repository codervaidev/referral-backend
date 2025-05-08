package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserGemRepo struct {
	DB *pgxpool.Pool
}

func NewUserGemRepo(db *pgxpool.Pool) *UserGemRepo {
	return &UserGemRepo{DB: db}
}

func (r *UserGemRepo) GetUserGems(ctx context.Context, userID uint) (int, error) {
	query := `
		SELECT points
		FROM referral_user
		WHERE user_id = $1
		LIMIT 1
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var points int
	for rows.Next() {
		err := rows.Scan(&points)
		if err != nil {
			return 0, err
		}
	}

	return points, nil
}

func (r *UserGemRepo) GetUserReferralCode(ctx context.Context, userID uint) (string, error) {
	query := `
		SELECT referer_code
		FROM referral_user
		WHERE user_id = $1
	`

	var referralCode string
	err := r.DB.QueryRow(ctx, query, userID).Scan(&referralCode)
	if err != nil {
		return "", err
	}

	return referralCode, nil
}

func (r *UserGemRepo) ValidateReferralCode(ctx context.Context, userID uint, referralCode string) (bool, error) {
	// Remove any surrounding quotes from the referral code
	referralCode = strings.Trim(referralCode, "'")
	
	
	
	query := `
		SELECT COUNT(*)
		FROM referral_user
		WHERE referer_code = $1
		AND user_id != $2
	`
	

	var count int
	err := r.DB.QueryRow(ctx, query, referralCode, userID).Scan(&count)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return false, err
	}

	return count > 0, nil
}
