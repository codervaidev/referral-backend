package repository

import (
	"context"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GemHistoryRepo struct {
	db *pgxpool.Pool
}


func NewGemHistoryRepo(db *pgxpool.Pool) *GemHistoryRepo {
	return &GemHistoryRepo{db: db}
}

func (r *GemHistoryRepo) GetGemHistory(ctx context.Context, userID uint) ([]models.GemHistoryResponse, error) {
	query := `
		SELECT gh.id, gh.amount, gh.user_id, gh.message, gh.type, gh.purchase_id, gh.created_at, u.phone, p.name, p."imageUrl" 
		FROM gems_history gh
		JOIN "User" u on u.id = user_id
		JOIN "Profile" p on p."userId" = u.id
		WHERE user_id = $1
		ORDER BY gh.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gemHistory []models.GemHistoryResponse
	for rows.Next() {
		var gh models.GemHistoryResponse
		err := rows.Scan(&gh.ID, &gh.Amount, &gh.UserID, &gh.Message, &gh.Type, &gh.PurchaseID, &gh.CreatedAt, &gh.Phone, &gh.Name, &gh.ImageUrl)
		if err != nil {
			return nil, err
		}
		
		// If this is a purchase type, get detailed purchase information
		if gh.Type == "purchase" && gh.PurchaseID != nil {
			purchaseDetails, err := r.getPurchaseDetails(ctx, *gh.PurchaseID)
			if err == nil {
				gh.PurchaseDetails = purchaseDetails
			}
		}
		
		gemHistory = append(gemHistory, gh)
	}

	return gemHistory, nil
}

// Add inserts a gem history record.
func (r *GemHistoryRepo) Add(ctx context.Context, userID uint, amount int, message, typ string, purchaseID *uuid.UUID) error {
	const q = `INSERT INTO gems_history (amount, user_id, message, type, purchase_id) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, amount, userID, message, typ, purchaseID)
	return err
}

// getPurchaseDetails retrieves detailed purchase information for a purchase history entry
func (r *GemHistoryRepo) getPurchaseDetails(ctx context.Context, purchaseID uuid.UUID) (*models.PurchaseHistoryDetails, error) {
	query := `
		SELECT 
			pr.id as purchase_id,
			dd.full_address as delivery_location,
			CASE 
				WHEN dd.district_id = 1 THEN 60 
				ELSE 120 
			END as delivery_fee,
			gh.amount as total_amount
		FROM purchase_referral pr
		JOIN delivery_details_referral dd ON dd.id = pr.delivery_id
		JOIN gems_history gh ON gh.purchase_id = pr.id
		WHERE pr.id = $1
	`
	
	var details models.PurchaseHistoryDetails
	err := r.db.QueryRow(ctx, query, purchaseID).Scan(
		&details.PurchaseID,
		&details.DeliveryLocation,
		&details.DeliveryFee,
		&details.TotalAmount,
	)
	if err != nil {
		return nil, err
	}
	
	// Get cart items for this purchase
	itemsQuery := `
		SELECT 
			ci.product_id,
			p.title as product_name,
			ci.variant_id,
			v.variant_name,
			ci.quantity,
			ci.price_at_add
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		LEFT JOIN variants v ON v.id = ci.variant_id
		JOIN carts c ON c.id = ci.cart_id
		JOIN purchase_referral pr ON pr.cart_id = c.id
		WHERE pr.id = $1
	`
	
	itemRows, err := r.db.Query(ctx, itemsQuery, purchaseID)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()
	
	for itemRows.Next() {
		var item models.PurchaseHistoryItem
		err := itemRows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.VariantID,
			&item.VariantName,
			&item.Quantity,
			&item.PriceAtAdd,
		)
		if err != nil {
			return nil, err
		}
		details.Items = append(details.Items, item)
	}
	
	return &details, nil
}
