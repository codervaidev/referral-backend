package repository

import (
	"context"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductPaymentRepo struct {
	DB *pgxpool.Pool
}

func NewProductPaymentRepo(db *pgxpool.Pool) *ProductPaymentRepo {
	return &ProductPaymentRepo{DB: db}
}

// GetReferralCodeStats returns statistics for a specific referral code
func (r *ProductPaymentRepo) GetReferralCodeStats(ctx context.Context, referralCode string) (*models.ReferralCodeStats, error) {
	query := `
		SELECT 
			referral_code,
			COUNT(CASE WHEN payment_status = 'success' THEN 1 END) as success_count,
			COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) as failed_count,
			COUNT(*) as total_count
		FROM product_payment 
		WHERE referral_code = $1
		GROUP BY referral_code
	`

	var stats models.ReferralCodeStats
	err := r.DB.QueryRow(ctx, query, referralCode).Scan(
		&stats.ReferralCode,
		&stats.SuccessCount,
		&stats.PendingCount,
		&stats.FailedCount,
		&stats.TotalCount,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAllReferralCodeStatsCount returns total number of distinct referral codes
func (r *ProductPaymentRepo) GetAllReferralCodeStatsCount(ctx context.Context) (int, error) {
    query := `
        SELECT COUNT(DISTINCT referral_code)
        FROM product_payment
        WHERE referral_code IS NOT NULL
    `

    var count int
    err := r.DB.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}

// GetAllReferralCodeStats returns statistics for all referral codes with pagination
func (r *ProductPaymentRepo) GetAllReferralCodeStats(ctx context.Context, limit, offset int) ([]models.ReferralCodeStats, error) {
    query := `
        SELECT 
            referral_code,
            COUNT(CASE WHEN payment_status = 'success' THEN 1 END) as success_count,
            COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_count,
            COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) as failed_count,
            COUNT(*) as total_count
        FROM product_payment 
        WHERE referral_code IS NOT NULL
        GROUP BY referral_code
        ORDER BY success_count DESC, total_count DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.DB.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stats []models.ReferralCodeStats
    for rows.Next() {
        var stat models.ReferralCodeStats
        err := rows.Scan(
            &stat.ReferralCode,
            &stat.SuccessCount,
            &stat.PendingCount,
            &stat.FailedCount,
            &stat.TotalCount,
        )
        if err != nil {
            return nil, err
        }
        stats = append(stats, stat)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return stats, nil
}

// GetSuccessfulPaymentsCount returns the total count of successful payments
func (r *ProductPaymentRepo) GetSuccessfulPaymentsCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM product_payment 
		WHERE payment_status = 'success' and referral_code is not null
	`

	var count int
	err := r.DB.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetSuccessfulPayments returns successful payments with pagination
func (r *ProductPaymentRepo) GetSuccessfulPayments(ctx context.Context, limit, offset int) ([]models.SuccessfulPayment, error) {
	query := `
		SELECT 
			amount_student_paid,
			invoice,
			user_id,
			referral_code,
			transaction_id,
			created_at
		FROM product_payment 
		WHERE payment_status = 'success' and referral_code is not null
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.SuccessfulPayment
	for rows.Next() {
		var payment models.SuccessfulPayment
		err := rows.Scan(
			&payment.AmountStudentPaid,
			&payment.Invoice,
			&payment.UserID,
			&payment.ReferralCode,
			&payment.TransactionID,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

// GetTotalPaymentStats returns overall payment statistics across all payments
func (r *ProductPaymentRepo) GetTotalPaymentStats(ctx context.Context) (*models.TotalPaymentStats, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN payment_status = 'success' THEN 1 END) as success_count,
			COUNT(CASE WHEN payment_status = 'pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) as failed_count,
			COUNT(*) as total_count
		FROM product_payment
		WHERE referral_code IS NOT NULL
	`

	var stats models.TotalPaymentStats
	err := r.DB.QueryRow(ctx, query).Scan(
		&stats.SuccessCount,
		&stats.PendingCount,
		&stats.FailedCount,
		&stats.TotalCount,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetTotalRevenueWithReferralCode calculates the total revenue from successful payments that have a referral code.
func (r *ProductPaymentRepo) GetTotalRevenueWithReferralCode(ctx context.Context) (float64, error) {
    query := `
        SELECT COALESCE(SUM(CAST(amount_student_paid AS DOUBLE PRECISION)), 0)
        FROM product_payment
        WHERE payment_status = 'success' AND referral_code IS NOT NULL
    `

    var totalRevenue float64
    err := r.DB.QueryRow(ctx, query).Scan(&totalRevenue)
    if err != nil {
        return 0, err
    }

    return totalRevenue, nil
} 