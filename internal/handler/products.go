package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/repository"
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
}

// GetAll returns every product in the catalogue.
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    products, err := h.Repo.GetAll(r.Context())
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

    product, err := h.Repo.GetByID(r.Context(), id)
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

    products, err := h.Repo.GetByCategoryID(r.Context(), catID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(products)
} 