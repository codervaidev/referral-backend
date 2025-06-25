package router

import (
	"github.com/codervaidev/referral-backend/internal/handler"
	"github.com/codervaidev/referral-backend/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(db *pgxpool.Pool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api").Subrouter()

	// Add middleware
	api.Use(middleware.CORSMiddleware)
	api.Use(middleware.MetricsMiddleware)

	handler := handler.New(db)

	// Health check
	api.HandleFunc("/health", handler.HealthCheck).Methods("GET")
	
	// Metrics endpoint
	api.Handle("/metrics", promhttp.Handler()).Methods("GET")
	
	// Gems routes
	handler.RegisterGemRoutes(api)
	handler.RegisterUserGemRoutes(api)
	handler.RegisterGemHistoryRoutes(api)
	// Register products routes
	handler.RegisterProductRoutes(api)
	// Categories routes
	handler.RegisterCategoryRoutes(api)
	
	// User details routes
	handler.RegisterUserDetailsRoutes(api)
	// Product wishlist routes
	handler.RegisterProductWishlistRoutes(api)
	// Cart routes
	handler.RegisterCartRoutes(api)

	return router
}
