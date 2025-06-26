package repository

import (
	"context"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GemHistoryRepo struct {
	db *pgxpool.Pool
}


func NewGemHistoryRepo(db *pgxpool.Pool) *GemHistoryRepo {
	return &GemHistoryRepo{db: db}
}

func (r *GemHistoryRepo) GetGemHistory(ctx context.Context, userID uint) ([]models.GemHistoryResponse, error) {
	query := `
		SELECT gh.id, gh.amount, gh.user_id, gh.message, gh.type,gh.created_at,u.phone,p.name,p."imageUrl" 
		FROM gems_history gh
		JOIN "User" u on u.id = user_id
		JOIN "Profile" p on p."userId" = u.id
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gemHistory []models.GemHistoryResponse
	for rows.Next() {
		var gh models.GemHistoryResponse
		err := rows.Scan(&gh.ID, &gh.Amount, &gh.UserID, &gh.Message, &gh.Type, &gh.CreatedAt, &gh.Phone, &gh.Name, &gh.ImageUrl)
		if err != nil {
			return nil, err
		}
		gemHistory = append(gemHistory, gh)
	}

	return gemHistory, nil
}

// Add inserts a gem history record.
func (r *GemHistoryRepo) Add(ctx context.Context, userID uint, amount int, message, typ string) error {
	const q = `INSERT INTO gems_history (amount, user_id, message, type) VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(ctx, q, amount, userID, message, typ)
	return err
}
