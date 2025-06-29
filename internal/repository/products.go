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
               (w.product_id IS NOT NULL) AS is_wishlisted,
               v.id as variant_id, v.variant_name, v.pics, v.variant_type
        FROM products p
        LEFT JOIN product_wishlist w ON (w.product_id = p.id AND w.user_id = $1 and w.variant_id is null) or (w.product_id = p.id AND w.user_id = $1 and w.variant_id is not null and w.variant_id = v.id)
        LEFT JOIN variants v ON v.product_id = p.id
        ORDER BY p.id, v.id`

    rows, err := r.DB.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    productMap := make(map[int]*models.Product)
    
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
            variantID *int
            variantName *string
            rawVariantPics []byte
            variantType *string
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted, &variantID, &variantName, &rawVariantPics, &variantType); err != nil {
            return nil, err
        }
        
        // Get or create product in map
        product, exists := productMap[p.ID]
        if !exists {
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
            p.Variants = []models.Variant{}
            productMap[p.ID] = &p
            product = &p
        }
        
        // Add variant if it exists
        if variantID != nil {
            variant := models.Variant{
                ID:          *variantID,
                ProductID:   &p.ID,
                VariantName: variantName,
                VariantType: variantType,
            }
            
            if len(rawVariantPics) > 0 {
                var pics []string
                if err := json.Unmarshal(rawVariantPics, &pics); err != nil {
                    return nil, err
                }
                variant.Pics = &pics
            }
            
            product.Variants = append(product.Variants, variant)
        }
    }
    
    // Convert map values to slice
    for _, product := range productMap {
        products = append(products, *product)
    }

    return products, nil
}

// GetByID fetches a single product row by its id, and whether it is wishlisted
// by the specified user.
func (r *ProductRepo) GetByID(ctx context.Context, id int, userID uint) (*models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted,
               v.id as variant_id, v.variant_name, v.pics, v.variant_type
        FROM products p
        LEFT JOIN product_wishlist w ON (w.product_id = p.id AND w.user_id = $2 and w.variant_id is null) or (w.product_id = p.id AND w.user_id = $2 and w.variant_id is not null and w.variant_id = v.id)
        LEFT JOIN variants v ON v.product_id = p.id
        WHERE p.id=$1
        ORDER BY v.id`

    rows, err := r.DB.Query(ctx, query, id, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var product *models.Product
    
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
            variantID *int
            variantName *string
            rawVariantPics []byte
            variantType *string
        )

        var recFor []int32
        err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted, &variantID, &variantName, &rawVariantPics, &variantType)
        if err != nil {
            return nil, err
        }
        
        if product == nil {
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
            p.Variants = []models.Variant{}
            product = &p
        }
        
        // Add variant if it exists
        if variantID != nil {
            variant := models.Variant{
                ID:          *variantID,
                ProductID:   &p.ID,
                VariantName: variantName,
                VariantType: variantType,
            }
            
            if len(rawVariantPics) > 0 {
                var pics []string
                if err := json.Unmarshal(rawVariantPics, &pics); err != nil {
                    return nil, err
                }
                variant.Pics = &pics
            }
            
            product.Variants = append(product.Variants, variant)
        }
    }
    
    if product == nil {
        return nil, err
    }
    
    return product, nil
}

func (r *ProductRepo) GetByCategoryID(ctx context.Context, categoryID int, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted,
               v.id as variant_id, v.variant_name, v.pics, v.variant_type
        FROM products p
        LEFT JOIN product_wishlist w ON (w.product_id = p.id AND w.user_id = $2 and w.variant_id is null) or (w.product_id = p.id AND w.user_id = $2 and w.variant_id is not null and w.variant_id = v.id)
        LEFT JOIN variants v ON v.product_id = p.id
        WHERE p.category_id=$1
        ORDER BY p.id, v.id`

    rows, err := r.DB.Query(ctx, query, categoryID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    productMap := make(map[int]*models.Product)
    
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
            variantID *int
            variantName *string
            rawVariantPics []byte
            variantType *string
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted, &variantID, &variantName, &rawVariantPics, &variantType); err != nil {
            return nil, err
        }
        
        // Get or create product in map
        product, exists := productMap[p.ID]
        if !exists {
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
            p.Variants = []models.Variant{}
            productMap[p.ID] = &p
            product = &p
        }
        
        // Add variant if it exists
        if variantID != nil {
            variant := models.Variant{
                ID:          *variantID,
                ProductID:   &p.ID,
                VariantName: variantName,
                VariantType: variantType,
            }
            
            if len(rawVariantPics) > 0 {
                var pics []string
                if err := json.Unmarshal(rawVariantPics, &pics); err != nil {
                    return nil, err
                }
                variant.Pics = &pics
            }
            
            product.Variants = append(product.Variants, variant)
        }
    }
    
    // Convert map values to slice
    for _, product := range productMap {
        products = append(products, *product)
    }

    return products, nil
}

// GetByRecommendedFor returns products that match the recommended_for class id.
func (r *ProductRepo) GetByRecommendedFor(ctx context.Context, classID int, userID uint) ([]models.Product, error) {
    const query = `
        SELECT p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating,
               p.recommended_for, p.image_urls, p.vendor,
               (w.product_id IS NOT NULL) AS is_wishlisted,
               v.id as variant_id, v.variant_name, v.pics, v.variant_type
        FROM products p
        LEFT JOIN product_wishlist w ON (w.product_id = p.id AND w.user_id = $2 and w.variant_id is null) or (w.product_id = p.id AND w.user_id = $2 and w.variant_id is not null and w.variant_id = v.id)
        LEFT JOIN variants v ON v.product_id = p.id
        WHERE $1 = ANY(p.recommended_for)
        ORDER BY p.id, v.id`

    rows, err := r.DB.Query(ctx, query, classID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []models.Product
    productMap := make(map[int]*models.Product)
    
    for rows.Next() {
        var (
            p          models.Product
            rawImgURLs []byte
            isWishlisted bool
            variantID *int
            variantName *string
            rawVariantPics []byte
            variantType *string
        )
        var recFor []int32
        if err := rows.Scan(&p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor, &isWishlisted, &variantID, &variantName, &rawVariantPics, &variantType); err != nil {
            return nil, err
        }
        
        // Get or create product in map
        product, exists := productMap[p.ID]
        if !exists {
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
            p.Variants = []models.Variant{}
            productMap[p.ID] = &p
            product = &p
        }
        
        // Add variant if it exists
        if variantID != nil {
            variant := models.Variant{
                ID:          *variantID,
                ProductID:   &p.ID,
                VariantName: variantName,
                VariantType: variantType,
            }
            
            if len(rawVariantPics) > 0 {
                var pics []string
                if err := json.Unmarshal(rawVariantPics, &pics); err != nil {
                    return nil, err
                }
                variant.Pics = &pics
            }
            
            product.Variants = append(product.Variants, variant)
        }
    }
    
    // Convert map values to slice
    for _, product := range productMap {
        products = append(products, *product)
    }

    return products, nil
} 