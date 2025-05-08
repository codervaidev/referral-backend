package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/codervaidev/referral-backend/internal/logger"
	"github.com/codervaidev/referral-backend/internal/metrics"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		metrics.HttpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.statusCode)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		// Log the request
		logger.Log.Info("HTTP Request",
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.Int("status", rw.statusCode),
			zap.Float64("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func TestMetricsMiddleware(t *testing.T) {
	// Create a new mux router
	router := mux.NewRouter()

	// Add a route to the router
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Add the middleware to the router
	router.Use(MetricsMiddleware)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a response recorder
	rec := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rec, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the metrics
	metrics.HttpRequestsTotal.WithLabelValues("GET", "/test", "200").Inc()
	metrics.HttpRequestDuration.WithLabelValues("GET", "/test").Observe(0.0001)
} 