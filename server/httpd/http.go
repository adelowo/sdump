package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func New(cfg config.Config,
	urlRepo sdump.URLRepository,
	ingestRepo sdump.IngestRepository,
	logger *logrus.Entry,
	sseServer *sse.Server,
) *http.Server {
	return &http.Server{
		Handler: buildRoutes(cfg, logger, urlRepo, ingestRepo, sseServer),
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
	}
}

func buildRoutes(cfg config.Config,
	logger *logrus.Entry,
	urlRepo sdump.URLRepository,
	ingestRepo sdump.IngestRepository,
	sseServer *sse.Server,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(writeRequestIDHeader)

	urlHandler := &urlHandler{
		cfg:        cfg,
		urlRepo:    urlRepo,
		logger:     logger,
		ingestRepo: ingestRepo,
		sseServer:  sseServer,
	}

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

func retrieveRequestID(r *http.Request) string { return middleware.GetReqID(r.Context()) }

var tracer = otel.Tracer("sdump")

func getTracer(ctx context.Context,
	r *http.Request, operationName string,
) (context.Context, trace.Span, string) {
	ctx, span := tracer.Start(ctx, operationName)

	rid := retrieveRequestID(r)

	span.SetAttributes(attribute.String("request_id", rid))

	return ctx, span, rid
}
