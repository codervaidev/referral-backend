package repository

import (
	"context"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
    DB *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
    return &CategoryRepo{DB: db}
}

// GetAll returns every category present in the `categories` table.
func (r *CategoryRepo) GetAll(ctx context.Context) ([]models.Category, error) {
    const query = `SELECT id, "name", description FROM categories`

    rows, err := r.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cats []models.Category
    for rows.Next() {
        var c models.Category
        if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
            return nil, err
        }
        cats = append(cats, c)
    }
    return cats, nil
} 