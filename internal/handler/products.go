package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/codervaidev/referral-backend/internal/config"
	"github.com/codervaidev/referral-backend/internal/middleware"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// ProductHandler wires HTTP requests to the underlying ProductRepo.
type ProductHandler struct {
	Repo *repository.ProductRepo
}

func (h *Handler) RegisterProductRoutes(r *mux.Router) {
	repo := repository.NewProductRepo(h.DB)
	ph := &ProductHandler{Repo: repo}

	// List all products
	r.HandleFunc("/products", ph.GetAll).Methods("GET", "OPTIONS")
	// Retrieve a single product by id
	r.HandleFunc("/products/{id}", ph.GetByID).Methods("GET", "OPTIONS")
	// Fetch products by category id
	r.HandleFunc("/products/category/{category_id}", ph.GetByCategoryID).Methods("GET", "OPTIONS")
	// Fetch products by recommended_for class
	r.HandleFunc("/products/recommended/{class_id}", ph.GetByRecommendedFor).Methods("GET", "OPTIONS")
}

// GetAll returns every product in the catalogue.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, _ := extractUserID(r)

	products, err := h.Repo.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// GetByID fetches a single product row based on the `id` path parameter.
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	userID, _ := extractUserID(r)

	product, err := h.Repo.GetByID(r.Context(), id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Fetch variants for this product
	vr := repository.NewVariantRepo(h.Repo.DB)
	variants, err := vr.GetByProductID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	product.Variants = variants

	json.NewEncoder(w).Encode(product)
}

// GetByCategoryID returns all products under a given category.
func (h *ProductHandler) GetByCategoryID(w http.ResponseWriter, r *http.Request) {
	catStr := mux.Vars(r)["category_id"]
	catID, err := strconv.Atoi(catStr)
	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	userID, _ := extractUserID(r)

	products, err := h.Repo.GetByCategoryID(r.Context(), catID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// GetByRecommendedFor fetches products where recommended_for matches the provided class id.
func (h *ProductHandler) GetByRecommendedFor(w http.ResponseWriter, r *http.Request) {
	clsStr := mux.Vars(r)["class_id"]
	clsID, err := strconv.Atoi(clsStr)
	if err != nil {
		http.Error(w, "Invalid class id", http.StatusBadRequest)
		return
	}

	userID, _ := extractUserID(r)

	products, err := h.Repo.GetByRecommendedFor(r.Context(), clsID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// extractUserID inspects the Authorization header for a Bearer token and, if
// present and valid, returns the user ID encoded in the JWT. The second return
// value indicates whether a valid user ID was extracted.
func extractUserID(r *http.Request) (uint, bool) {
	// First, attempt to fetch user ID set by JWT middleware if present.
	if val := r.Context().Value(middleware.UserIDKey); val != nil {
		if s, ok := val.(string); ok {
			if uid64, err := strconv.ParseUint(s, 10, 32); err == nil {
				return uint(uid64), true
			}
		}
	}

	// Fallback: parse Authorization header manually (supports public routes without middleware).
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, false
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	cfg := config.Load()

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, false
	}

	return uint(userIDFloat), true
}
