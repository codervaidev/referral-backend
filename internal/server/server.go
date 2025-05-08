package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/codervaidev/referral-backend/internal/config"
	"github.com/codervaidev/referral-backend/internal/db"
	"github.com/codervaidev/referral-backend/internal/router"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
}

func New(cfg *config.Config) *Server {
	pgPool, err := db.NewPostgres(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to DB: %v", err))
	}

	r := router.New(pgPool)

	return &Server{
		db: pgPool,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: r,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.db != nil {
		s.db.Close()
	}
	return s.httpServer.Shutdown(ctx)
}
