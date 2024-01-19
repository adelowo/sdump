package httpd

import (
	"fmt"
	"net/http"

	"github.com/adelowo/sdump/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(cfg config.Config) *http.Server {
	return &http.Server{
		Handler: buildRoutes(),
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
	}
}

func buildRoutes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)

	return router
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", r.Context().Value(middleware.RequestIDKey).(string))
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }
