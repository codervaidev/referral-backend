package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductPayment struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	ProductID           uuid.UUID  `json:"product_id" db:"product_id"`
	PaymentMethod       *string    `json:"payment_method" db:"payment_method"`
	AmountStudentPaid   *string    `json:"amount_student_paid" db:"amount_student_paid"`
	AmountReceived      *string    `json:"amount_received" db:"amount_received"`
	PaymentStatus       string     `json:"payment_status" db:"payment_status"`
	PaymentDate         *string    `json:"payment_date" db:"payment_date"`
	TransactionID       string     `json:"transaction_id" db:"transaction_id"`
	GatewayResponse     *string    `json:"gateway_response" db:"gateway_response"`
	CreatedAt           *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at" db:"updated_at"`
	UserID              int        `json:"user_id" db:"user_id"`
	Invoice             *string    `json:"invoice" db:"invoice"`
	Source              *string    `json:"source" db:"source"`
	DeliveryID          *uuid.UUID `json:"delivery_id" db:"delivery_id"`
	UtmID               *uuid.UUID `json:"utm_id" db:"utm_id"`
	IpnResponse         *string    `json:"ipn_response" db:"ipn_response"`
	IsDelivered         *bool      `json:"is_delivered" db:"is_delivered"`
	Comment             *string    `json:"comment" db:"comment"`
	BkashTrxID          *string    `json:"bkash_trxid" db:"bkash_trxid"`
	BkashPaymentID      *string    `json:"bkash_paymentid" db:"bkash_paymentid"`
	BkashResponse       *string    `json:"bkash_response" db:"bkash_response"`
	ReferralCode        *string    `json:"referral_code" db:"referral_code"`
}

type ReferralCodeStats struct {
	ReferralCode   string `json:"referral_code"`
	SuccessCount   int    `json:"success_count"`
	PendingCount   int    `json:"pending_count"`
	FailedCount    int    `json:"failed_count"`
	TotalCount     int    `json:"total_count"`
}

type SuccessfulPayment struct {
	AmountStudentPaid *string `json:"amount_student_paid" db:"amount_student_paid"`
	Invoice           *string `json:"invoice" db:"invoice"`
	UserID            int     `json:"user_id" db:"user_id"`
	ReferralCode      *string `json:"referral_code" db:"referral_code"`
	TransactionID     *string `json:"transaction_id" db:"transaction_id"`
	CreatedAt         *time.Time `json:"created_at" db:"created_at"`
}

type TotalPaymentStats struct {
	SuccessCount int `json:"success_count"`
	PendingCount int `json:"pending_count"`
	FailedCount  int `json:"failed_count"`
	TotalCount   int `json:"total_count"`
}

type PaginatedSuccessfulPayments struct {
	Data       []SuccessfulPayment `json:"data"`
	TotalCount int                 `json:"total_count"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type PaginatedReferralCodeStats struct {
    Data       []ReferralCodeStats `json:"data"`
    TotalCount int                 `json:"total_count"`
    Page       int                 `json:"page"`
    PageSize   int                 `json:"page_size"`
    TotalPages int                 `json:"total_pages"`
}

// NEW_STRUCT
// TotalRevenue represents the aggregated revenue from successful payments with a referral code.
type TotalRevenue struct {
    TotalRevenue float64 `json:"total_revenue"`
}