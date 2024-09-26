package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/datastore/postgres"
	"github.com/adelowo/sdump/datastore/sqlite"
	"github.com/adelowo/sdump/server/httpd"
	"github.com/r3labs/sse/v2"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

func getDatabase(cfg *config.Config) (*bun.DB, error) {
	if cfg.HTTP.Database.Driver == "sqllite" {
		return sqlite.New(cfg.HTTP.Database.DSN, cfg.HTTP.Database.LogQueries)
	}

	return postgres.New(cfg.HTTP.Database.DSN, cfg.HTTP.Database.LogQueries)
}

func createHTTPCommand(cmd *cobra.Command, cfg *config.Config) {
	command := &cobra.Command{
		Use:   "http",
		Short: "Start/run the HTTP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			sig := make(chan os.Signal, 1)

			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			cleanupFn := func(_ context.Context) error {
				return nil
			}

			var err error

			if cfg.HTTP.OTEL.IsEnabled {
				cleanupFn, err = initTracer(cfg)
				if err != nil {
					return err
				}
			}

			lvl, err := logrus.ParseLevel(cfg.LogLevel)
			if err != nil {
				return err
			}

			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(lvl)

			var tokensPerMinute uint64 = 60
			if cfg.HTTP.RateLimit.RequestsPerMinute > 0 {
				tokensPerMinute = cfg.HTTP.RateLimit.RequestsPerMinute
			}

			ratelimitStore, err := memorystore.New(&memorystore.Config{
				Tokens:   tokensPerMinute,
				Interval: time.Minute,
			})
			if err != nil {
				return err
			}

			// TODO: use appropriate stores
			// Or just return the stores themselves, not the inner db
			db, err := getDatabase(cfg)
			if err != nil {
				return err
			}
			urlStore := postgres.NewURLRepositoryTable(db)
			ingestStore := postgres.NewIngestRepository(db)
			userStore := postgres.NewUserRepositoryTable(db)

			hostName, err := os.Hostname()
			if err != nil {
				return err
			}

			logger := logrus.WithField("host", hostName).
				WithField("module", "http.server")

			sseServer := sse.New()

			httpServer := httpd.New(*cfg, urlStore, ingestStore,
				userStore, logger, sseServer, ratelimitStore)

			go func() {
				logger.Debug("starting HTTP server")
				if err := httpServer.ListenAndServe(); err != nil {
					logger.WithError(err).Fatal("could not start http server")
				}
			}()

			<-sig
			if err := cleanupFn(context.Background()); err != nil {
				logger.WithError(err).Error("could not properly shut down OTEL")
			}

			if err := httpServer.Shutdown(context.Background()); err != nil {
				logger.WithError(err).Error("could not shut down http server")
			}

			if err := db.Close(); err != nil {
				logger.WithError(err).Error("could not shut down database connection")
			}

			sseServer.Close()
			ratelimitStore.Close(context.Background())

			return nil
		},
	}

	cmd.AddCommand(command)
}

func initTracer(cfg *config.Config) (func(context.Context) error, error) {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if !cfg.HTTP.OTEL.UseTLS {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.HTTP.OTEL.Endpoint),
		),
	)
	if err != nil {
		return nil, err
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.HTTP.OTEL.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	return exporter.Shutdown, nil
}
