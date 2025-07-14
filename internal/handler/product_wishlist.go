package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/config"
	"github.com/codervaidev/referral-backend/internal/middleware"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

// ProductWishlistHandler wires HTTP requests to ProductWishlist repository.
type ProductWishlistHandler struct {
	Repo *repository.ProductWishlistRepo
}

// RegisterProductWishlistRoutes registers wishlist routes under /wishlist prefix.
func (h *Handler) RegisterProductWishlistRoutes(r *mux.Router) {
	repo := repository.NewProductWishlistRepo(h.DB)
	wh := &ProductWishlistHandler{Repo: repo}

	cfg := config.Load()
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	// Fetch wishlist products
	r.Handle("/wishlist", jwtMiddleware.Middleware(http.HandlerFunc(wh.GetWishlist))).Methods("GET", "OPTIONS")
	// Add product to wishlist
	r.Handle("/wishlist", jwtMiddleware.Middleware(http.HandlerFunc(wh.AddProduct))).Methods("POST", "OPTIONS")
	// Remove product from wishlist
	r.Handle("/wishlist/{product_id}", jwtMiddleware.Middleware(http.HandlerFunc(wh.RemoveProduct))).Methods("DELETE", "OPTIONS")
	// Get wishlist count
	r.Handle("/wishlist/count", jwtMiddleware.Middleware(http.HandlerFunc(wh.GetWishlistCount))).Methods("GET", "OPTIONS")
}

// AddProduct adds a product (and optional variant) to the authenticated user's wishlist.
// Expects JSON body: { "product_id": 123 } and optional query param ?variant_id=456.
func (h *ProductWishlistHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		ProductID int `json:"product_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse optional variant_id query parameter.
	variantID, err := getVariantIDQueryParam(r)
	if err != nil {
		http.Error(w, "invalid variant_id", http.StatusBadRequest)
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Add(r.Context(), uint(userIDUint64), payload.ProductID, variantID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "added"})
}

// RemoveProduct deletes a product from the user's wishlist. The product_id is
// supplied as a path param.
func (h *ProductWishlistHandler) RemoveProduct(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pidStr := mux.Vars(r)["product_id"]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	// Parse optional variant_id query parameter for deletion.
	variantID, err := getVariantIDQueryParam(r)
	if err != nil {
		http.Error(w, "invalid variant_id", http.StatusBadRequest)
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Remove(r.Context(), uint(userIDUint64), pid, variantID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
}

// GetWishlist returns the list of products present in the authenticated user's wishlist.
func (h *ProductWishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	products, err := h.Repo.GetProductsByUserID(r.Context(), uint(userIDUint64))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Optionally fetch variants for each product (similar to product detail handler)
	// vr := repository.NewVariantRepo(h.Repo.DB)
	for i := range products {
		// variants, err := vr.GetByProductID(r.Context(), products[i].ID)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// Product is wishlisted only when variant_id is NULL (i.e., list entry applied to whole product)
		products[i].IsWishlisted = true
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductWishlistHandler) GetWishlistCount(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	userIDUint64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	count, err := h.Repo.GetWishlistCount(r.Context(), uint(userIDUint64))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
