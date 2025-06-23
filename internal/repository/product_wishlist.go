package repository

import (
	"context"
	"encoding/json"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductWishlistRepo handles CRUD operations for the `product_wishlist` table.
type ProductWishlistRepo struct {
    DB *pgxpool.Pool
}

func NewProductWishlistRepo(db *pgxpool.Pool) *ProductWishlistRepo {
    return &ProductWishlistRepo{DB: db}
}

// Add inserts a row into product_wishlist. An error is returned if the row
// already exists (i.e. the user already added the product to their wishlist).
func (r *ProductWishlistRepo) Add(ctx context.Context, userID uint, productID int) error {
    const query = `
        INSERT INTO product_wishlist (user_id, product_id)
        VALUES ($1, $2) ON CONFLICT DO NOTHING`

    _, err := r.DB.Exec(ctx, query, userID, productID)
    return err
}

// Remove deletes a wishlist entry for the given user and product.
func (r *ProductWishlistRepo) Remove(ctx context.Context, userID uint, productID int) error {
    const query = `DELETE FROM product_wishlist WHERE user_id=$1 AND product_id=$2`
    _, err := r.DB.Exec(ctx, query, userID, productID)
    return err
}

// GetProductsByUserID returns the list of products that the user has added to
// their wishlist. The function performs a join to fetch full product details.
// Variants are not loaded here – callers can fetch them separately if needed.
func (r *ProductWishlistRepo) GetProductsByUserID(ctx context.Context, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating, p.recommended_for, p.image_urls, p.vendor
        FROM products p
        INNER JOIN product_wishlist w ON w.product_id = p.id
        WHERE w.user_id = $1`

    rows, err := r.DB.Query(ctx, query, userID)
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
        var recFor []int32

        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor); err != nil {
            return nil, err
        }

        if len(recFor) > 0 {
            ints := make([]int, len(recFor))
            for i, v := range recFor {
                ints[i] = int(v)
            }
            p.RecommendedFor = &ints
        }

        if len(rawImgURLs) > 0 {
            var urls []string
            if err := json.Unmarshal(rawImgURLs, &urls); err != nil {
                return nil, err
            }
            p.ImageURLs = &urls
        }

        // Since these products are fetched from the user's wishlist,
        // they are, by definition, wish-listed.
        p.IsWishlisted = true

        products = append(products, p)
    }

    return products, nil
} 