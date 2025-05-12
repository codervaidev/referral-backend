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

func (r *GemHistoryRepo) GetGemHistory(ctx context.Context, userID uint) ([]models.GemHistory, error) {
	query := `
		SELECT id, amount, user_id, message, type
		FROM gems_history
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gemHistory []models.GemHistory
	for rows.Next() {
		var gh models.GemHistory
		err := rows.Scan(&gh.ID, &gh.Amount, &gh.UserID, &gh.Message, &gh.Type)
		if err != nil {
			return nil, err
		}
		gemHistory = append(gemHistory, gh)
	}

	return gemHistory, nil
}
