package handler

import (
	"encoding/json"
	"net/http"

	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

// CalculatorHandler exposes endpoints that return gem pricing for study packages.
// All data is read-only and currently static.
// It intentionally does not require authentication because pricing is public.

type CalculatorHandler struct {
    Repo *repository.CalculatorRepo
}

// RegisterCalculatorRoutes mounts the /calculator route on the provided router.
func (h *Handler) RegisterCalculatorRoutes(r *mux.Router) {
    repo := repository.NewCalculatorRepo(nil)
    ch := &CalculatorHandler{Repo: repo}

    r.HandleFunc("/calculator", ch.GetPricing).Methods("GET", "OPTIONS")
}

// GetPricing returns the full pricing table as JSON.
func (h *CalculatorHandler) GetPricing(w http.ResponseWriter, r *http.Request) {
    data := h.Repo.GetPricing()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
} 