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
    rows, err := r.DB.Query(ctx, "SELECT * FROM gems_store")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var gems []models.Gem
    for rows.Next() {
        var g models.Gem
        err := rows.Scan(&g.ID, &g.Name, &g.Image, &g.GemsCount, &g.IsActive, &g.Title, &g.Description, &g.Category)
        if err != nil {
            return nil, err
        }
        gems = append(gems, g)
    }
    return gems, nil
}

func (r *GemRepo) GetByID(ctx context.Context, id string) (*models.Gem, error) {
    row := r.DB.QueryRow(ctx, "SELECT id, name, image, gems_count, is_active, title, description, category FROM gems_store WHERE id=$1", id)
    var g models.Gem
    err := row.Scan(&g.ID, &g.Name, &g.Image, &g.GemsCount, &g.IsActive, &g.Title, &g.Description, &g.Category)
    if err != nil {
        return nil, err
    }
    return &g, nil
}

func (r *GemRepo) Create(ctx context.Context, g models.Gem) (uuid.UUID, error) {
    var id uuid.UUID
    row := r.DB.QueryRow(ctx, `INSERT INTO gems_store
        (name, image, gems_count, is_active, title, description, category, link, variant, picture_1, picture_2, picture_3, picture_4)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id`,
        g.Name, g.Image, g.GemsCount, g.IsActive, g.Title, g.Description, g.Category, g.Link, g.Variant, g.Picture1, g.Picture2, g.Picture3, g.Picture4)
    err := row.Scan(&id)
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}

func (r *GemRepo) Update(ctx context.Context, g models.Gem) error {
    _, err := r.DB.Exec(ctx, "UPDATE gems_store SET name=$1, image=$2, gems_count=$3, is_active=$4, title=$5, description=$6, category=$7 WHERE id=$8",
                    g.Name, g.Image, g.GemsCount, g.IsActive, g.Title, g.Description, g.Category, g.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *GemRepo) Delete(ctx context.Context, id string) error {
    _, err := r.DB.Exec(ctx, "DELETE FROM gems_store WHERE id=$1", id)
    return err
}

func (r *GemRepo) GetGemsCountWithUserID(ctx context.Context, userID string) (int, error) {
    var gemsCount int
    err := r.DB.QueryRow(ctx, "SELECT gems_count FROM gems_store WHERE id=$1", userID).Scan(&gemsCount)
    if err != nil {
        return 0, err
    }
    return gemsCount, nil
}

func (r *GemRepo) GetLeaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
    query := `
        SELECT 
            u.id AS user_id,
            p.name AS student_name,
            p."imageUrl" AS student_image,
            pc.class,
            r.points AS total_gems
        FROM referral_user r
        LEFT JOIN "User" u ON r.user_id = u.id
        LEFT JOIN "Profile" p ON p."userId" = u.id
        LEFT JOIN profile_classes pc ON pc.profile_id = p.id
        ORDER BY r.points DESC
        LIMIT 10
    `
    
    rows, err := r.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var entries []models.LeaderboardEntry
    for rows.Next() {
        var entry models.LeaderboardEntry
        err := rows.Scan(&entry.UserID, &entry.StudentName, &entry.StudentImage, &entry.Class, &entry.TotalGems)
        if err != nil {
            return nil, err
        }
        entries = append(entries, entry)
    }
    
    return entries, nil
}
