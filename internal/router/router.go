package router

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"app_aggregator/internal/domain"
	"app_aggregator/internal/handlers"
	"app_aggregator/internal/middleware"
	"app_aggregator/internal/ratelimit"
)

type HTTPServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewHTTPServer(
	organizationService domain.OrganizationService,
	loanApplicationService domain.LoanApplicationService,
	logger *slog.Logger,
) *HTTPServer {
	mux := http.NewServeMux()

	organizationHandler := handlers.NewHTTPOrganizationHandler(organizationService, logger)
	loanApplicationHandler := handlers.NewHTTPLoanApplicationHandler(loanApplicationService, logger)

	registerRoutes(mux, organizationHandler, loanApplicationHandler)

	rateLimitConfig := &ratelimit.Config{
		RequestsPerMinute: 100,
		WindowSize:        1 * time.Minute,
		CleanupInterval:   1 * time.Minute,
		BlockDuration:     1 * time.Minute,
	}

	handler := middleware.Chain(
		mux,
		middleware.Logger(logger),
		middleware.RateLimitWithLogger(rateLimitConfig, logger),
		middleware.CORS(),
		middleware.Recovery(logger),
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		server: server,
		logger: logger,
	}
}

func (s *HTTPServer) Start() error {
	s.logger.Info("Starting HTTP server", slog.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

func registerRoutes(
	mux *http.ServeMux,
	orgHandler *handlers.HTTPOrganizationHandler,
	loanHandler *handlers.HTTPLoanApplicationHandler,
) {
	mux.HandleFunc("GET /api/v1/organizations", orgHandler.GetAll)
	mux.HandleFunc("GET /api/v1/organizations/{uuid}", orgHandler.GetByID)

	mux.HandleFunc("POST /api/v1/admin/organizations", orgHandler.Create)
	mux.HandleFunc("PATCH /api/v1/admin/organizations/{uuid}", orgHandler.Update)
	mux.HandleFunc("DELETE /api/v1/admin/organizations/{uuid}", orgHandler.Delete)

	mux.HandleFunc("GET /api/v1/loan_applications", loanHandler.GetAll)
	mux.HandleFunc("GET /api/v1/loan_applications/{uuid}", loanHandler.GetByID)
	mux.HandleFunc("POST /api/v1/loan_applications", loanHandler.Create)
	mux.HandleFunc("PATCH /api/v1/loan_applications/{uuid}", loanHandler.Update)
	mux.HandleFunc("DELETE /api/v1/loan_applications/{uuid}", loanHandler.Delete)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
}
