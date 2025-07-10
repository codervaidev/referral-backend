package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codervaidev/referral-backend/internal/config"
	"github.com/codervaidev/referral-backend/internal/middleware"

	"github.com/codervaidev/referral-backend/internal/models"
	"github.com/codervaidev/referral-backend/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GemHandler struct {
	Repo *repository.GemRepo
}

func (h *Handler) RegisterGemRoutes(r *mux.Router) {
	repo := repository.NewGemRepo(h.DB)
	gh := &GemHandler{Repo: repo}

	// JWT middleware for routes requiring authentication
	cfg := config.Load()
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	r.HandleFunc("/gems", gh.GetAll).Methods("GET")
	r.HandleFunc("/gems/leaderboard", gh.GetLeaderboard).Methods("GET", "OPTIONS")
	r.Handle("/gems/leaderboard/position", jwtMiddleware.Middleware(http.HandlerFunc(gh.GetUserLeaderboardPosition))).Methods("GET", "OPTIONS")
	r.HandleFunc("/gems/{id}", gh.GetByID).Methods("GET")
	r.HandleFunc("/gems", gh.Create).Methods("POST")
	r.HandleFunc("/gems/{id}", gh.Update).Methods("PUT")
	r.HandleFunc("/gems/{id}", gh.Delete).Methods("DELETE")
}

func (h *GemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	gems, err := h.Repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(gems)
}

func (h *GemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	gem, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(gem)
}

func (h *GemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var gem models.Gem
	if err := json.NewDecoder(r.Body).Decode(&gem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.Repo.Create(r.Context(), gem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	gem.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(gem)
}

func (h *GemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var gem models.Gem
	if err := json.NewDecoder(r.Body).Decode(&gem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gem.ID, _ = uuid.Parse(id)

	if err := h.Repo.Update(r.Context(), gem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *GemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GemHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := h.Repo.GetLeaderboard(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// GetUserLeaderboardPosition returns the authenticated user's rank on the leaderboard.
func (h *GemHandler) GetUserLeaderboardPosition(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uid64, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	position, err := h.Repo.GetUserLeaderboardPosition(r.Context(), uint(uid64))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"position": position})
}
