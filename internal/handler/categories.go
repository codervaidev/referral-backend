package handler

import (
	"encoding/json"
	"net/http"

	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

// CategoryHandler exposes endpoints for category resources.
type CategoryHandler struct {
	Repo *repository.CategoryRepo
}

func (h *Handler) RegisterCategoryRoutes(r *mux.Router) {
	repo := repository.NewCategoryRepo(h.DB)
	ch := &CategoryHandler{Repo: repo}

	r.HandleFunc("/categories", ch.GetAll).Methods("GET", "OPTIONS")
}

// GetAll returns the list of categories.
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	cats, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cats)
}
