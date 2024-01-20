package httpd

import (
	"fmt"
	"net/http"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func New(cfg config.Config,
	urlRepo sdump.URLRepository,
	logger *logrus.Entry,
) *http.Server {
	return &http.Server{
		Handler: buildRoutes(cfg, logger, urlRepo),
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
	}
}

func buildRoutes(cfg config.Config,
	logger *logrus.Entry,
	urlRepo sdump.URLRepository,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)

	urlHandler := &urlHandler{
		cfg:     cfg,
		urlRepo: urlRepo,
		logger:  logger,
	}

	router.Post("/", urlHandler.create)
	router.Post("/{reference}", urlHandler.ingest)

	return router
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", r.Context().Value(middleware.RequestIDKey).(string))
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }
