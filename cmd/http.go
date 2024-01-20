package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/datastore/postgres"
	"github.com/adelowo/sdump/server/httpd"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func createHTTPCommand(cmd *cobra.Command, cfg *config.Config) {
	command := &cobra.Command{
		Use:  "http",
		Long: "Start/run the HTTP server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sig := make(chan os.Signal, 1)

			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			if cfg.Log == "" {
				cfg.Log = "error"
			}

			lvl, err := logrus.ParseLevel(cfg.Log)
			if err != nil {
				return err
			}

			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(lvl)

			db, err := postgres.New(cfg.HTTP.Database.DSN, cfg.HTTP.Database.LogQueries)
			if err != nil {
				return err
			}

			urlStore := postgres.NewURLRepositoryTable(db)
			ingestStore := postgres.NewIngestRepository(db)

			hostName, err := os.Hostname()
			if err != nil {
				return err
			}

			logger := logrus.WithField("host", hostName).
				WithField("module", "http.server")

			sseServer := sse.New()

			httpServer := httpd.New(*cfg, urlStore, ingestStore,
				logger, sseServer)

			go func() {
				logger.Debug("starting HTTP server")
				if err := httpServer.ListenAndServe(); err != nil {
					logger.WithError(err).Fatal("could not start http server")
				}
			}()

			<-sig

			if err := db.Close(); err != nil {
				logger.WithError(err).Error("could not shut down database connection")
			}

			if err := httpServer.Shutdown(context.Background()); err != nil {
				logger.WithError(err).Error("could not shut down http server")
			}

			sseServer.Close()

			return nil
		},
	}

	cmd.AddCommand(command)
}
