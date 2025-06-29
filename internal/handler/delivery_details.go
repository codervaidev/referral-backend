package handler

import (
	"encoding/json"
	"net/http"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/gorilla/mux"
)

type DeliveryDetailsHandler struct {
	Repo *repository.DeliveryDetailsRepo
}

func (h *Handler) RegisterDeliveryDetailsRoutes(r *mux.Router) {
	repo := repository.NewDeliveryDetailsRepo(h.DB)
	ddh := &DeliveryDetailsHandler{Repo: repo}

	r.HandleFunc("/delivery-details", ddh.Create).Methods("POST", "OPTIONS")
}

// Create saves delivery details for the authenticated user.
func (h *DeliveryDetailsHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := extractUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.DeliveryDetail
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payload.UserID = int(userID)

	id, err := h.Repo.Create(r.Context(), payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Delivery details saved",
		"id":      id,
	})
}
