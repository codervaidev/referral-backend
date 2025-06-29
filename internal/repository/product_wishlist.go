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

// Add inserts a row into product_wishlist. If variantID is nil, the column is left NULL.
// An error is returned if the row already exists (i.e. the user already added the product/variant combination).
func (r *ProductWishlistRepo) Add(ctx context.Context, userID uint, productID int, variantID *int) error {
    var err error
    if variantID != nil {
        const query = `
            INSERT INTO product_wishlist (user_id, product_id, variant_id)
            VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
        _, err = r.DB.Exec(ctx, query, userID, productID, *variantID)
    } else {
        const query = `
            INSERT INTO product_wishlist (user_id, product_id)
            VALUES ($1, $2) ON CONFLICT DO NOTHING`
        _, err = r.DB.Exec(ctx, query, userID, productID)
    }
    return err
}

// Remove deletes a wishlist entry for the given user, product, and optional variant.
func (r *ProductWishlistRepo) Remove(ctx context.Context, userID uint, productID int, variantID *int) error {
    var err error
    if variantID != nil {
        const query = `DELETE FROM product_wishlist WHERE user_id=$1 AND product_id=$2 AND variant_id=$3`
        _, err = r.DB.Exec(ctx, query, userID, productID, *variantID)
    } else {
        const query = `DELETE FROM product_wishlist WHERE user_id=$1 AND product_id=$2 AND variant_id IS NULL`
        _, err = r.DB.Exec(ctx, query, userID, productID)
    }
    return err
}

// GetProductsByUserID returns the list of products that the user has added to
// their wishlist. The function performs a join to fetch full product details.
// Variants are not loaded here – callers can fetch them separately if needed.
func (r *ProductWishlistRepo) GetProductsByUserID(ctx context.Context, userID uint) ([]models.Product, error) {
    const query = `
        SELECT 
  p.id, p.category_id, p.link, p.title, p.description, 
  p.price, p.stock, p.sold, p.wishlist_count, p.rating, 
  p.recommended_for, p.image_urls, p.vendor,
  TRUE AS is_wishlisted,
  COALESCE(
    json_agg(
      json_build_object(
        'id', v.id,
        'product_id', v.product_id,
        'variant_name', v.variant_name,
        'pics', v.pics,
        'variant_type', v.variant_type
      )
    ) FILTER (WHERE v.id IS NOT NULL),
    '[]'::json
  ) AS variants
FROM product_wishlist w
JOIN products p ON p.id = w.product_id
LEFT JOIN variants v ON v.id = w.variant_id
WHERE w.user_id = $1
GROUP BY 
  p.id, p.category_id, p.link, p.title, p.description, 
  p.price, p.stock, p.sold, p.wishlist_count, p.rating, 
  p.recommended_for, p.image_urls, p.vendor;
`

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
            rawVariants []byte
        )
        var recFor []int32

        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &p.IsWishlisted, &rawVariants); err != nil {
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

        // Parse variants from JSON
        if len(rawVariants) > 0 {
            var variants []models.Variant
            if err := json.Unmarshal(rawVariants, &variants); err != nil {
                return nil, err
            }
            p.Variants = variants
        }

        products = append(products, p)
    }

    return products, nil
} 