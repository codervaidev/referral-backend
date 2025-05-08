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

type UserGemHandler struct {
	Repo *repository.UserGemRepo
}

func (h *Handler) RegisterUserGemRoutes(r *mux.Router) {
	userGemRepo := repository.NewUserGemRepo(h.DB)
	userGemHandler := &UserGemHandler{Repo: userGemRepo}
	cfg := config.Load()
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	r.Handle("/user-gems", jwtMiddleware.Middleware(http.HandlerFunc(userGemHandler.GetUserGems))).Methods("GET")
	r.Handle("/user-referral-code", jwtMiddleware.Middleware(http.HandlerFunc(userGemHandler.GetUserReferralCode))).Methods("GET")
	r.Handle("/validate-referral-code", jwtMiddleware.Middleware(http.HandlerFunc(userGemHandler.ValidateReferralCode))).Methods("GET")
}

func (h *UserGemHandler) GetUserGems(w http.ResponseWriter, r *http.Request) {
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

	gems, err := h.Repo.GetUserGems(r.Context(), uint(userIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(gems)
}

func (h *UserGemHandler) GetUserReferralCode(w http.ResponseWriter, r *http.Request) {
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

	referralCode, err := h.Repo.GetUserReferralCode(r.Context(), uint(userIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(referralCode)
}

func (h *UserGemHandler) ValidateReferralCode(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	referralCode := r.URL.Query().Get("referral_code")
	if referralCode == "" {
		http.Error(w, "Referral code is required", http.StatusBadRequest)
		return
	}

	userIDInt, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	valid, err := h.Repo.ValidateReferralCode(r.Context(), uint(userIDInt), referralCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(valid)
}
