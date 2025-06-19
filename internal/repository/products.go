package repository

import (
	"context"
	"encoding/json"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
    DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
    return &ProductRepo{DB: db}
}

// GetAll retrieves every product from the products table.
func (r *ProductRepo) GetAll(ctx context.Context) ([]models.Product, error) {
    const query = `
        SELECT id, category_id, link, title, description, price, stock, sold, wishlist_count, rating, image_urls, vendor
        FROM products`

    rows, err := r.DB.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
        )
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &rawImgURLs, &p.Vendor); err != nil {
            return nil, err
        }
        // Decode the JSONB image_urls column.
        if len(rawImgURLs) > 0 {
            var urls []string
            if err := json.Unmarshal(rawImgURLs, &urls); err != nil {
                return nil, err
            }
            p.ImageURLs = &urls
        }
        products = append(products, p)
    }

    return products, nil
}

// GetByID fetches a single product row by its id.
func (r *ProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
    const query = `
        SELECT id, category_id, link, title, description, price, stock, sold, wishlist_count, rating, image_urls, vendor
        FROM products WHERE id=$1`

    var (
        p          models.Product
        rawImgURLs []byte
    )

    err := r.DB.QueryRow(ctx, query, id).Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &rawImgURLs, &p.Vendor)
    if err != nil {
        return nil, err
    }
    if len(rawImgURLs) > 0 {
        var urls []string
        if err := json.Unmarshal(rawImgURLs, &urls); err != nil {
            return nil, err
        }
        p.ImageURLs = &urls
    }
    return &p, nil
}

func (r *ProductRepo) GetByCategoryID(ctx context.Context, categoryID int) ([]models.Product, error) {
    const query = `
        SELECT id, category_id, link, title, description, price, stock, sold, wishlist_count, rating, image_urls, vendor
        FROM products WHERE category_id=$1`

    rows, err := r.DB.Query(ctx, query, categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
        )
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &rawImgURLs, &p.Vendor); err != nil {
            return nil, err
        }
        if len(rawImgURLs) > 0 {
            var urls []string
            if err := json.Unmarshal(rawImgURLs, &urls); err != nil {
                return nil, err
            }
            p.ImageURLs = &urls
        }
        products = append(products, p)
    }

    return products, nil
} 