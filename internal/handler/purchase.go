package handler

import (
	"encoding/json"
	"net/http"

	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PurchaseHandler manages purchase flow.
type PurchaseHandler struct {
	CartRepo           *repository.CartRepo
	GemRepo            *repository.UserGemRepo
	GemHistoryRepo     *repository.GemHistoryRepo
	PurchaseRepo       *repository.PurchaseRepo
	DeliveryDetailsRepo *repository.DeliveryDetailsRepo
}

func (h *Handler) RegisterPurchaseRoutes(r *mux.Router) {
	ph := &PurchaseHandler{
		CartRepo:           repository.NewCartRepo(h.DB),
		GemRepo:            repository.NewUserGemRepo(h.DB),
		GemHistoryRepo:     repository.NewGemHistoryRepo(h.DB),
		PurchaseRepo:       repository.NewPurchaseRepo(h.DB),
		DeliveryDetailsRepo: repository.NewDeliveryDetailsRepo(h.DB),
	}

	r.HandleFunc("/purchase", ph.Create).Methods("POST", "OPTIONS")
}

// purchase request body
type purchaseReq struct {
	DeliveryID uuid.UUID `json:"delivery_id"`
}

func (ph *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req purchaseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fetch active cart (create if none) – must have items
	cart, err := ph.CartRepo.GetActiveCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if cart == nil || len(cart.Items) == 0 {
		http.Error(w, "cart is empty", http.StatusBadRequest)
		return
	}

	// compute total
	total := 0.0
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.PriceAtAdd
	}
	pointsRequired := int(total)

	// Get delivery details to determine district-based pricing
	deliveryDetails, err := ph.DeliveryDetailsRepo.GetByID(r.Context(), req.DeliveryID)
	if err != nil {
		http.Error(w, "Invalid delivery details", http.StatusBadRequest)
		return
	}

	// Add delivery fee based on district
	deliveryFee := 0
	if deliveryDetails.DistrictID == 1 {
		// Inside Dhaka - add 60 gems
		deliveryFee = 60
	} else {
		// Outside Dhaka - add 120 gems
		deliveryFee = 120
	}

	// Add delivery fee to total points required
	pointsRequired += deliveryFee

	// deduct gems
	if err := ph.GemRepo.DeductPoints(r.Context(), userID, pointsRequired); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// add gem history for product purchase
	_ = ph.GemHistoryRepo.Add(r.Context(), userID, -(pointsRequired-deliveryFee), "Purchase deduction", "debit")
	
	// add gem history for delivery fee
	_ = ph.GemHistoryRepo.Add(r.Context(), userID, -deliveryFee, "Delivery fee", "debit")

	// mark cart checked_out
	if err := ph.CartRepo.Checkout(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create purchase
	purchaseID, err := ph.PurchaseRepo.Create(r.Context(), userID, req.DeliveryID, cart.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":           "Purchase successful",
		"purchase_id":       purchaseID,
		"gems_deducted":     pointsRequired,
		"product_cost":      pointsRequired - deliveryFee,
		"delivery_fee":      deliveryFee,
		"district_id":       deliveryDetails.DistrictID,
		"delivery_location": deliveryDetails.FullAddress,
	})
}
