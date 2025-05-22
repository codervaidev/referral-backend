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

type GemHistoryHandler struct {
	Repo *repository.GemHistoryRepo
}

func (h *Handler) RegisterGemHistoryRoutes(r *mux.Router) {
	gemHistoryRepo := repository.NewGemHistoryRepo(h.DB)
	gemHistoryHandler := &GemHistoryHandler{Repo: gemHistoryRepo}
	cfg := config.Load()
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	r.Handle("/gem-history", jwtMiddleware.Middleware(http.HandlerFunc(gemHistoryHandler.GetGemHistory))).Methods("GET","OPTIONS")
}

func (h *GemHistoryHandler) GetGemHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	gemHistory, err := h.Repo.GetGemHistory(r.Context(), uint(userIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gemHistory)
}