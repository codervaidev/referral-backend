package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserDetailsHandler struct {
	Repo *repository.UserDetailsRepo
}

func (h *Handler) RegisterUserDetailsRoutes(r *mux.Router) {
	repo := repository.NewUserDetailsRepo(h.DB)
	udh := &UserDetailsHandler{Repo: repo}

	r.HandleFunc("/user-details", udh.Create).Methods("POST", "OPTIONS")
	r.HandleFunc("/user-details/{user_id}", udh.GetByUserID).Methods("GET", "OPTIONS")
	r.HandleFunc("/user-details/{user_id}", udh.Update).Methods("PUT", "OPTIONS")
}

func (h *UserDetailsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var details models.UserDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Always generate a new UUID, ignore any id from the request
	details.ID = uuid.New()

	if err := h.Repo.Create(r.Context(), details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(details)
}

func (h *UserDetailsHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	details, err := h.Repo.GetByUserID(r.Context(), uint(userIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(details)
}

func (h *UserDetailsHandler) Update(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["user_id"]
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	var details models.UserDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	details.UserID = userIDInt
	if err := h.Repo.Update(r.Context(), details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(details)
}
