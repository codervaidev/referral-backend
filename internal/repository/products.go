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

// GetAll retrieves every product and marks them as wishlisted for the provided
// user ID, if any. Pass userID = 0 to indicate an unauthenticated request.
func (r *ProductRepo) GetAll(ctx context.Context, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted
        FROM products p
        LEFT JOIN product_wishlist w ON w.product_id = p.id AND w.user_id = $1`

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
            isWishlisted bool
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted); err != nil {
            return nil, err
        }
        if len(recFor) > 0 {
            ints := make([]int, len(recFor))
            for i, v := range recFor {
                ints[i] = int(v)
            }
            p.RecommendedFor = &ints
        }
        // Decode the JSONB image_urls column.
        if len(rawImgURLs) > 0 {
            var urls []string
            if err := json.Unmarshal(rawImgURLs, &urls); err != nil {
                return nil, err
            }
            p.ImageURLs = &urls
        }
        p.IsWishlisted = isWishlisted
        products = append(products, p)
    }

    return products, nil
}

// GetByID fetches a single product row by its id, and whether it is wishlisted
// by the specified user.
func (r *ProductRepo) GetByID(ctx context.Context, id int, userID uint) (*models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted
        FROM products p
        LEFT JOIN product_wishlist w ON w.product_id = p.id AND w.user_id = $2
        WHERE p.id=$1`

    var (
        p          models.Product
        rawImgURLs []byte
        isWishlisted bool
    )

    var recFor []int32
    err := r.DB.QueryRow(ctx, query, id, userID).Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted)
    if err != nil {
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
    p.IsWishlisted = isWishlisted
    return &p, nil
}

func (r *ProductRepo) GetByCategoryID(ctx context.Context, categoryID int, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted
        FROM products p
        LEFT JOIN product_wishlist w ON w.product_id = p.id AND w.user_id = $2
        WHERE p.category_id=$1`

    rows, err := r.DB.Query(ctx, query, categoryID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted); err != nil {
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
        p.IsWishlisted = isWishlisted
        products = append(products, p)
    }

    return products, nil
}

// GetByRecommendedFor returns products that match the recommended_for class id.
func (r *ProductRepo) GetByRecommendedFor(ctx context.Context, classID int, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted
        FROM products p
        LEFT JOIN product_wishlist w ON w.product_id = p.id AND w.user_id = $2
        WHERE $1 = ANY(p.recommended_for)`

    rows, err := r.DB.Query(ctx, query, classID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted); err != nil {
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
        p.IsWishlisted = isWishlisted
        products = append(products, p)
    }

    return products, nil
} 