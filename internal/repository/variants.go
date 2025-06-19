package repository

import (
	"context"
	"encoding/json"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantRepo struct {
    DB *pgxpool.Pool
}

func NewVariantRepo(db *pgxpool.Pool) *VariantRepo {
    return &VariantRepo{DB: db}
}

// GetByProductID returns all variants that belong to a single product.
func (r *VariantRepo) GetByProductID(ctx context.Context, productID int) ([]models.Variant, error) {
    const query = `SELECT id, product_id, variant_name, pics, variant_type FROM variants WHERE product_id=$1`

    rows, err := r.DB.Query(ctx, query, productID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var variants []models.Variant
    for rows.Next() {
        var (
            v       models.Variant
            rawPics []byte
        )
        if err := rows.Scan(&v.ID, &v.ProductID, &v.VariantName, &rawPics, &v.VariantType); err != nil {
            return nil, err
        }
        if len(rawPics) > 0 {
            var pics []string
            if err := json.Unmarshal(rawPics, &pics); err != nil {
                return nil, err
            }
            v.Pics = &pics
        }
        variants = append(variants, v)
    }

    return variants, nil
} 