package handler

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

// ProductPaymentHandler handles product payment related endpoints
type ProductPaymentHandler struct {
	Repo *repository.ProductPaymentRepo
}

// RegisterProductPaymentRoutes registers product payment routes
func (h *Handler) RegisterProductPaymentRoutes(r *mux.Router) {
	repo := repository.NewProductPaymentRepo(h.DB)
	pph := &ProductPaymentHandler{Repo: repo}

	// Get statistics for a specific referral code
	r.HandleFunc("/product-payment/referral-stats/{referral_code}", pph.GetReferralCodeStats).Methods("GET", "OPTIONS")
	// Get statistics for all referral codes
	r.HandleFunc("/product-payment/referral-stats", pph.GetAllReferralCodeStats).Methods("GET", "OPTIONS")
	// Get successful payments with specific fields
	r.HandleFunc("/product-payment/successful-payments", pph.GetSuccessfulPayments).Methods("GET", "OPTIONS")
	// Get total payment statistics
	r.HandleFunc("/product-payment/total-stats", pph.GetTotalPaymentStats).Methods("GET", "OPTIONS")
	// NEW_ROUTE: Get total revenue from successful payments with referral code
	r.HandleFunc("/product-payment/total-revenue", pph.GetTotalRevenue).Methods("GET", "OPTIONS")
}

// GetReferralCodeStats returns payment statistics for a specific referral code
func (h *ProductPaymentHandler) GetReferralCodeStats(w http.ResponseWriter, r *http.Request) {
	referralCode := mux.Vars(r)["referral_code"]
	if referralCode == "" {
		http.Error(w, "Referral code is required", http.StatusBadRequest)
		return
	}

	stats, err := h.Repo.GetReferralCodeStats(r.Context(), referralCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetAllReferralCodeStats returns payment statistics for all referral codes
func (h *ProductPaymentHandler) GetAllReferralCodeStats(w http.ResponseWriter, r *http.Request) {
    // Parse pagination parameters
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }
    pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
    if err != nil || pageSize < 1 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize

    // total distinct referral codes
    totalCount, err := h.Repo.GetAllReferralCodeStatsCount(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    stats, err := h.Repo.GetAllReferralCodeStats(r.Context(), pageSize, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

    response := models.PaginatedReferralCodeStats{
        Data:       stats,
        TotalCount: totalCount,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetSuccessfulPayments returns all successful payments with specific fields
func (h *ProductPaymentHandler) GetSuccessfulPayments(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	totalCount, err := h.Repo.GetSuccessfulPaymentsCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get paginated payments
	payments, err := h.Repo.GetSuccessfulPayments(r.Context(), pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Create paginated response
	response := models.PaginatedSuccessfulPayments{
		Data:       payments,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTotalPaymentStats returns overall payment statistics across all payments
func (h *ProductPaymentHandler) GetTotalPaymentStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Repo.GetTotalPaymentStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// NEW_HANDLER
// GetTotalRevenue returns the total revenue generated from successful payments that have a referral code.
func (h *ProductPaymentHandler) GetTotalRevenue(w http.ResponseWriter, r *http.Request) {
    revenue, err := h.Repo.GetTotalRevenueWithReferralCode(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.TotalRevenue{TotalRevenue: revenue})
}
