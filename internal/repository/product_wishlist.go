package repository

import (
	"context"
	"encoding/json"
	"fmt"

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
        fmt.Println("variantID value:", *variantID)
    } else {
        fmt.Println("variantID is nil")
    }
    fmt.Println("userID", userID)
    fmt.Println("productID", productID)
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
            products.id,
            products.category_id,
            products.link,
            products.title,
            products.description,
            products.price,
            products.stock,
            products.sold,
            products.wishlist_count,
            products.rating,
            products.recommended_for,
            products.image_urls,
            products.vendor,
            variants.id,
            variants.variant_name,
            CASE WHEN product_wishlist.variant_id IS NULL THEN true ELSE false END as is_wishlisted
        FROM product_wishlist
        JOIN products ON products.id = product_wishlist.product_id
        LEFT JOIN variants ON variants.id = product_wishlist.variant_id
        WHERE product_wishlist.user_id = $1
        `
    
    rows, err := r.DB.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    //fmt.Println("rows", rows)
    var products []models.Product
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            variantID *int
            variantName *string
            isWishlisted bool
        )
        var recFor []int32

        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &variantID, &variantName, &isWishlisted); err != nil {
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

        // store wishlist variant id for later use in handler
        p.WishlistVariantID = variantID

        // Set is_wishlisted field
        p.IsWishlisted = isWishlisted

        // If variantID is present, append to Variants array
        if variantID != nil {
            v := models.Variant{
                ID: *variantID,
                ProductID: &p.ID,
                VariantName: variantName,
                IsWishlisted: true,
            }
            p.Variants = append(p.Variants, v)
        }

        products = append(products, p)
        fmt.Println("products", p);
    }

    fmt.Printf("products array: %+v\n", products)

    // Or for pretty JSON:
    if b, err := json.MarshalIndent(products, "", "  "); err == nil {
        fmt.Println("products array (json):", string(b))
    }

    return products, nil
} 