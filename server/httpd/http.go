package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/telemetry"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/r3labs/sse/v2"
	"github.com/riandyrn/otelchi"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var createdURLMetrics = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sdump_created_urls",
	Help: "The number of urls that have been created in total",
})

var ingestedHTTPRequestsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sdump_ingested_http_requests",
	Help: "Total number of ingested HTTP requests",
})

var failedIngestedHTTPRequestsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "sdump_ingested_http_requests_failure",
	Help: "Total number of failed attempts to ingest HTTP requests",
})

func New(cfg config.Config,
	urlRepo sdump.URLRepository,
	ingestRepo sdump.IngestRepository,
	userRepo sdump.UserRepository,
	logger *logrus.Entry,
	sseServer *sse.Server,
) *http.Server {
	return &http.Server{
		Handler: buildRoutes(cfg, logger, urlRepo, ingestRepo,
			userRepo, sseServer),
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
	}
}

func buildRoutes(cfg config.Config,
	logger *logrus.Entry,
	urlRepo sdump.URLRepository,
	ingestRepo sdump.IngestRepository,
	userRepo sdump.UserRepository,
	sseServer *sse.Server,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)
	router.Use(jsonResponse)

	urlHandler := &urlHandler{
		cfg:        cfg,
		urlRepo:    urlRepo,
		logger:     logger,
		ingestRepo: ingestRepo,
		userRepo:   userRepo,
		sseServer:  sseServer,
	}

	router.Use(writeRequestIDHeader)

	if cfg.HTTP.Prometheus.IsEnabled {
		router.Use(telemetry.Collector(telemetry.Config{
			Username: cfg.HTTP.Prometheus.Username,
			Password: cfg.HTTP.Prometheus.Password,
		}, []string{"/"}))

		_ = prometheus.Register(ingestedHTTPRequestsCounter)
		_ = prometheus.Register(failedIngestedHTTPRequestsCounter)
		_ = prometheus.Register(createdURLMetrics)
	}

	router.Use(otelchi.Middleware("http-router", otelchi.WithChiRoutes(router)))

	router.Post("/", urlHandler.create)
	router.Post("/{reference}", urlHandler.ingest)
	router.Get("/events", sseServer.ServeHTTP)

	return router
}

func writeRequestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", r.Context().Value(middleware.RequestIDKey).(string))
		next.ServeHTTP(w, r)
	})
}

func jsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }

var tracer = otel.Tracer("sdump.http")

func getTracer(ctx context.Context,
	r *http.Request, operationName string,
) (context.Context, trace.Span, string) {
	ctx, span := tracer.Start(ctx, operationName)

	rid := retrieveRequestID(r)

	span.SetAttributes(attribute.String("request_id", rid))

	return ctx, span, rid
}
