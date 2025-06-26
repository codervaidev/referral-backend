package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CartRepo provides access to cart and cart_items tables.
type CartRepo struct {
    DB *pgxpool.Pool
}

func NewCartRepo(db *pgxpool.Pool) *CartRepo { return &CartRepo{DB: db} }

// ensureActiveCart returns the active cart id for the given user, creating it
// if necessary.
func (r *CartRepo) ensureActiveCart(ctx context.Context, userID uint) (uuid.UUID, error) {
    const selectQuery = `SELECT id FROM carts WHERE user_id=$1 AND status='active' LIMIT 1`
    var cartID uuid.UUID
    err := r.DB.QueryRow(ctx, selectQuery, userID).Scan(&cartID)
    if err == nil {
        return cartID, nil
    }
    if !errors.Is(err, pgx.ErrNoRows) {
        return uuid.Nil, err
    }

    const insertQuery = `INSERT INTO carts (user_id, status) VALUES ($1, 'active') RETURNING id`
    if err := r.DB.QueryRow(ctx, insertQuery, userID).Scan(&cartID); err != nil {
        return uuid.Nil, err
    }
    return cartID, nil
}

// GetActiveCart loads the active cart and its items for the user. If no active
// cart exists, an empty cart is returned with Items set to an empty slice.
func (r *CartRepo) GetActiveCart(ctx context.Context, userID uint) (*models.Cart, error) {
    const cartQuery = `SELECT id, user_id, status, created_at, updated_at FROM carts WHERE user_id=$1 AND status='active' LIMIT 1`
    var c models.Cart
    if err := r.DB.QueryRow(ctx, cartQuery, userID).Scan(&c.ID, &c.UserID, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            // No cart yet – return empty with Items=nil
            return nil, nil
        }
        return nil, err
    }

    items, err := r.getItemsByCartID(ctx, c.ID)
    if err != nil {
        return nil, err
    }
    c.Items = items
    return &c, nil
}

func (r *CartRepo) getItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]models.CartItem, error) {
    const query = `
        SELECT ci.id, ci.cart_id, ci.product_id, ci.variant_id, ci.quantity, ci.price_at_add, ci.created_at, ci.updated_at,
               p.id, p.category_id, p.link, p.title, p.description, p.price, p.stock, p.sold, p.wishlist_count, p.rating, p.recommended_for, p.image_urls, p.vendor,
               v.id, v.product_id, v.variant_name, v.pics, v.variant_type
        FROM cart_items ci
        JOIN products p ON p.id = ci.product_id
        LEFT JOIN variants v ON v.id = ci.variant_id
        WHERE ci.cart_id=$1`

    rows, err := r.DB.Query(ctx, query, cartID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.CartItem
    for rows.Next() {
        var (
            it models.CartItem
            p  models.Product
            rawImgURLs []byte
            recFor []int32
            rawPics []byte
            vID *int
            varProdID *int
            varName *string
            varType *string
        )
        if err := rows.Scan(&it.ID, &it.CartID, &it.ProductID, &it.VariantID, &it.Quantity, &it.PriceAtAdd, &it.CreatedAt, &it.UpdatedAt,
            &p.ID, &p.CategoryID, &p.Link, &p.Title, &p.Description, &p.Price, &p.Stock, &p.Sold, &p.WishlistCount, &p.Rating, &recFor, &rawImgURLs, &p.Vendor,
            &vID, &varProdID, &varName, &rawPics, &varType); err != nil {
            return nil, err
        }
        if len(rawImgURLs) > 0 {
            if err := p.ScanImageURLs(rawImgURLs); err != nil {
                return nil, err
            }
        }
        if len(recFor) > 0 {
            ints := make([]int, len(recFor))
            for i, v := range recFor {
                ints[i] = int(v)
            }
            p.RecommendedFor = &ints
        }
        it.Product = &p

        if vID != nil {
            var variant models.Variant
            if vID != nil {
                variant.ID = *vID
            }
            variant.ProductID = varProdID
            variant.VariantName = varName
            variant.VariantType = varType
            // pics
            if len(rawPics) > 0 {
                var pics []string
                if err := json.Unmarshal(rawPics, &pics); err != nil {
                    return nil, err
                }
                variant.Pics = &pics
            }
            it.Variant = &variant
        }
        items = append(items, it)
    }

    return items, nil
}

// AddItem adds quantity units of a product (and optional variant) to the
// user's active cart, creating the cart if it does not exist. If the item
// already exists in the cart, its quantity is increased.
func (r *CartRepo) AddItem(ctx context.Context, userID uint, productID int, variantID *int, quantity int) error {
    if quantity <= 0 {
        return errors.New("quantity must be greater than zero")
    }

    cartID, err := r.ensureActiveCart(ctx, userID)
    if err != nil {
        return err
    }

    // Fetch current product price for historical accuracy.
    const priceQuery = `SELECT price FROM products WHERE id=$1`
    var priceAtAdd float64
    if err := r.DB.QueryRow(ctx, priceQuery, productID).Scan(&priceAtAdd); err != nil {
        return err
    }

    // Attempt to update existing row first
    var result pgconn.CommandTag
    if variantID != nil {
        const upd = `UPDATE cart_items SET quantity = quantity + $1, updated_at = now() WHERE cart_id=$2 AND product_id=$3 AND variant_id=$4`
        result, err = r.DB.Exec(ctx, upd, quantity, cartID, productID, *variantID)
    } else {
        const upd = `UPDATE cart_items SET quantity = quantity + $1, updated_at = now() WHERE cart_id=$2 AND product_id=$3 AND variant_id IS NULL`
        result, err = r.DB.Exec(ctx, upd, quantity, cartID, productID)
    }
    if err != nil {
        return err
    }
    if result.RowsAffected() > 0 {
        return nil // updated
    }

    // Insert new row
    if variantID != nil {
        const ins = `INSERT INTO cart_items (cart_id, product_id, variant_id, quantity, price_at_add) VALUES ($1,$2,$3,$4,$5)`
        _, err = r.DB.Exec(ctx, ins, cartID, productID, *variantID, quantity, priceAtAdd)
    } else {
        const ins = `INSERT INTO cart_items (cart_id, product_id, quantity, price_at_add) VALUES ($1,$2,$3,$4)`
        _, err = r.DB.Exec(ctx, ins, cartID, productID, quantity, priceAtAdd)
    }
    return err
}

// UpdateItemQuantity sets the quantity of an existing item. If quantity==0 the
// item is removed.
func (r *CartRepo) UpdateItemQuantity(ctx context.Context, userID uint, productID int, variantID *int, quantity int) error {
    cartID, err := r.ensureActiveCart(ctx, userID)
    if err != nil {
        return err
    }

    if quantity <= 0 {
        return r.RemoveItem(ctx, userID, productID, variantID)
    }

    var cmd pgconn.CommandTag
    if variantID != nil {
        const upd = `UPDATE cart_items SET quantity=$1, updated_at=now() WHERE cart_id=$2 AND product_id=$3 AND variant_id=$4`
        cmd, err = r.DB.Exec(ctx, upd, quantity, cartID, productID, *variantID)
    } else {
        const upd = `UPDATE cart_items SET quantity=$1, updated_at=now() WHERE cart_id=$2 AND product_id=$3 AND variant_id IS NULL`
        cmd, err = r.DB.Exec(ctx, upd, quantity, cartID, productID)
    }
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("item not found in cart")
    }
    return nil
}

// RemoveItem deletes an item from the cart.
func (r *CartRepo) RemoveItem(ctx context.Context, userID uint, productID int, variantID *int) error {
    cartID, err := r.ensureActiveCart(ctx, userID)
    if err != nil {
        return err
    }

    var cmd pgconn.CommandTag
    if variantID != nil {
        const del = `DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2 AND variant_id=$3`
        cmd, err = r.DB.Exec(ctx, del, cartID, productID, *variantID)
    } else {
        const del = `DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2 AND variant_id IS NULL`
        cmd, err = r.DB.Exec(ctx, del, cartID, productID)
    }
    if err != nil {
        return err
    }
    if cmd.RowsAffected() == 0 {
        return errors.New("item not found in cart")
    }
    return nil
}

// Clear marks the current active cart as "abandoned". The rows remain for
// historical purposes. A new active cart will be lazily created on the next
// AddItem call.
func (r *CartRepo) Clear(ctx context.Context, userID uint) error {
    const upd = `UPDATE carts SET status='abandoned', updated_at=now() WHERE user_id=$1 AND status='active'`
    _, err := r.DB.Exec(ctx, upd, userID)
    return err
}

// Checkout finalises the active cart by flipping its status and leaving the
// current items intact. A real implementation would create an order, etc.
func (r *CartRepo) Checkout(ctx context.Context, userID uint) error {
    const upd = `UPDATE carts SET status='checked_out', updated_at=now() WHERE user_id=$1 AND status='active'`
    _, err := r.DB.Exec(ctx, upd, userID)
    return err
}

// GetLatestCheckedOutCart returns the most recently checked-out cart for the user.
func (r *CartRepo) GetLatestCheckedOutCart(ctx context.Context, userID uint) (*models.Cart, error) {
    const cartQuery = `SELECT id, user_id, status, created_at, updated_at FROM carts WHERE user_id=$1 AND status='checked_out' ORDER BY updated_at DESC LIMIT 1`
    var c models.Cart
    if err := r.DB.QueryRow(ctx, cartQuery, userID).Scan(&c.ID, &c.UserID, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }

    items, err := r.getItemsByCartID(ctx, c.ID)
    if err != nil {
        return nil, err
    }
    c.Items = items
    return &c, nil
}

// GetActiveCartTotal returns the summed price (quantity * price_at_add) of all
// items in the user's active cart. If no cart or items exist, total is 0.
func (r *CartRepo) GetActiveCartTotal(ctx context.Context, userID uint) (float64, error) {
    const q = `SELECT COALESCE(SUM(ci.quantity * ci.price_at_add), 0)
               FROM carts c
               JOIN cart_items ci ON ci.cart_id = c.id
               WHERE c.user_id=$1 AND c.status='active'`

    var total float64
    if err := r.DB.QueryRow(ctx, q, userID).Scan(&total); err != nil {
        return 0, err
    }
    return total, nil
} 