package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

// CartHandler handles /cart endpoints.
type CartHandler struct {
    Repo *repository.CartRepo
}

// RegisterCartRoutes adds the cart endpoints to the router.
func (h *Handler) RegisterCartRoutes(r *mux.Router) {
    cr := repository.NewCartRepo(h.DB)
    ch := &CartHandler{Repo: cr}

    // Always require JWT middleware – these routes operate on the current user.
    // The top-level router already attaches middleware; ensuring user id extraction inside handlers.

    r.HandleFunc("/cart", ch.GetCart).Methods("GET", "OPTIONS")
    r.HandleFunc("/cart/items", ch.AddItem).Methods("POST", "OPTIONS")
    r.HandleFunc("/cart/items/{product_id}", ch.UpdateItemQuantity).Methods("PUT", "OPTIONS")
    r.HandleFunc("/cart/items/{product_id}", ch.RemoveItem).Methods("DELETE", "OPTIONS")
    r.HandleFunc("/cart/items", ch.Clear).Methods("DELETE", "OPTIONS")
    r.HandleFunc("/cart/checkout", ch.Checkout).Methods("POST", "OPTIONS")
}

// GetCart returns the authenticated user's active cart including items.
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    cart, err := h.Repo.GetActiveCart(r.Context(), userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if cart == nil {
        // Return empty structure for consistency
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(struct {
            Items []interface{} `json:"items"`
        }{Items: []interface{}{}})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cart)
}

// AddItem expects JSON body: {"product_id":123, "variant_id":456, "quantity":2}
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req struct {
        ProductID int  `json:"product_id"`
        VariantID *int `json:"variant_id,omitempty"`
        Quantity  int  `json:"quantity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.Repo.AddItem(r.Context(), userID, req.ProductID, req.VariantID, req.Quantity); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]any{
        "message": "Item added to cart",
    })
}

// UpdateItemQuantity updates quantity for product_id path parameter. Optional variant_id query param.
func (h *CartHandler) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    prodStr := mux.Vars(r)["product_id"]
    productID, err := strconv.Atoi(prodStr)
    if err != nil {
        http.Error(w, "invalid product_id", http.StatusBadRequest)
        return
    }

    variantID, err := getVariantIDQueryParam(r)
    if err != nil {
        http.Error(w, "invalid variant_id", http.StatusBadRequest)
        return
    }

    var req struct {
        Quantity int `json:"quantity"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.Repo.UpdateItemQuantity(r.Context(), userID, productID, variantID, req.Quantity); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]any{
        "message": "Cart item updated",
    })
}

// RemoveItem deletes an item from the cart.
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    prodStr := mux.Vars(r)["product_id"]
    productID, err := strconv.Atoi(prodStr)
    if err != nil {
        http.Error(w, "invalid product_id", http.StatusBadRequest)
        return
    }

    variantID, err := getVariantIDQueryParam(r)
    if err != nil {
        http.Error(w, "invalid variant_id", http.StatusBadRequest)
        return
    }

    if err := h.Repo.RemoveItem(r.Context(), userID, productID, variantID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]any{
        "message": "Item removed from cart",
    })
}

// Clear wipes all items from the cart.
func (h *CartHandler) Clear(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if err := h.Repo.Clear(r.Context(), userID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]any{
        "message": "Cart cleared",
    })
}

// Checkout finalises the cart.
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
    userID, ok := extractUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if err := h.Repo.Checkout(r.Context(), userID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]any{
        "message": "Checkout successful",
    })
}

// getVariantIDQueryParam parses the optional variant_id query parameter.
func getVariantIDQueryParam(r *http.Request) (*int, error) {
    variantStr := r.URL.Query().Get("variant_id")
    if variantStr == "" {
        return nil, nil
    }
    v, err := strconv.Atoi(variantStr)
    if err != nil {
        return nil, err
    }
    return &v, nil
} 