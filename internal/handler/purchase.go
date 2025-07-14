package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

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
	ProductRepo        *repository.ProductRepo
}

func (h *Handler) RegisterPurchaseRoutes(r *mux.Router) {
	ph := &PurchaseHandler{
		CartRepo:           repository.NewCartRepo(h.DB),
		GemRepo:            repository.NewUserGemRepo(h.DB),
		GemHistoryRepo:     repository.NewGemHistoryRepo(h.DB),
		PurchaseRepo:       repository.NewPurchaseRepo(h.DB),
		DeliveryDetailsRepo: repository.NewDeliveryDetailsRepo(h.DB),
		ProductRepo:        repository.NewProductRepo(h.DB),
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

	// validate stock availability for all items
	for _, item := range cart.Items {
		product, err := ph.ProductRepo.GetByID(r.Context(), item.ProductID, userID)
		if err != nil {
			http.Error(w, "Error checking product stock", http.StatusInternalServerError)
			return
		}
		if product.Stock == nil || *product.Stock < item.Quantity {
			http.Error(w, fmt.Sprintf("Insufficient stock for product %d. Available: %d, Requested: %d", item.ProductID, *product.Stock, item.Quantity), http.StatusBadRequest)
			return
		}
	}

	// deduct gems
	if err := ph.GemRepo.DeductPoints(r.Context(), userID, pointsRequired); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decrease stock for each product in cart
	for _, item := range cart.Items {
		if err := ph.ProductRepo.DecreaseStockAndUpdateSold(r.Context(), item.ProductID, item.Quantity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

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

	// add single gems_history record for the entire purchase
	_ = ph.GemHistoryRepo.Add(r.Context(), userID, -pointsRequired, "তুমি " + strconv.Itoa(pointsRequired) + " পয়েন্ট ব্যবহার করে নিচের পণ্যগুলো ক্রয় করেছো", "purchase", &purchaseID)

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
