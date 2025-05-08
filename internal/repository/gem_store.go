package repository

import (
	"context"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GemRepo struct {
    DB *pgxpool.Pool
}

func NewGemRepo(db *pgxpool.Pool) *GemRepo {
    return &GemRepo{DB: db}
}

func (r *GemRepo) GetAll(ctx context.Context) ([]models.Gem, error) {
    rows, err := r.DB.Query(ctx, "SELECT id, name, image, gems_count, is_active FROM gems_store")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var gems []models.Gem
    for rows.Next() {
        var g models.Gem
        err := rows.Scan(&g.ID, &g.Name, &g.Image, &g.GemsCount, &g.IsActive)
        if err != nil {
            return nil, err
        }
        gems = append(gems, g)
    }
    return gems, nil
}

func (r *GemRepo) GetByID(ctx context.Context, id string) (*models.Gem, error) {
    row := r.DB.QueryRow(ctx, "SELECT id, name, image, gems_count, is_active FROM gems_store WHERE id=$1", id)
    var g models.Gem
    err := row.Scan(&g.ID, &g.Name, &g.Image, &g.GemsCount, &g.IsActive)
    if err != nil {
        return nil, err
    }
    return &g, nil
}

func (r *GemRepo) Create(ctx context.Context, g models.Gem) (uuid.UUID, error)	 {
    var id uuid.UUID
    row := r.DB.QueryRow(ctx, "INSERT INTO gems_store(name, image, gems_count, is_active) VALUES($1, $2, $3, $4) returning id",
        g.Name, g.Image, g.GemsCount, g.IsActive)
    err := row.Scan(&id)
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}

func (r *GemRepo) Update(ctx context.Context, g models.Gem) error {
    _, err := r.DB.Exec(ctx, "UPDATE gems_store SET name=$1, image=$2, gems_count=$3, is_active=$4 WHERE id=$5",
        g.Name, g.Image, g.GemsCount, g.IsActive, g.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *GemRepo) Delete(ctx context.Context, id string) error {
    _, err := r.DB.Exec(ctx, "DELETE FROM gems_store WHERE id=$1", id)
    return err
}
