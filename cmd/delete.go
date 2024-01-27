package main

import (
	"context"
	"time"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/datastore/postgres"
	"github.com/spf13/cobra"
)

func createDeleteCommand(rootCmd *cobra.Command, cfg *config.Config) {
	cmd := &cobra.Command{
		Use:     "delete-http",
		Aliases: []string{"d"},
		Short:   "Deletes all old HTTP requests to preserve DB space",
		RunE: func(_ *cobra.Command, _ []string) error {
			db, err := postgres.New(cfg.HTTP.Database.DSN, cfg.HTTP.Database.LogQueries)
			if err != nil {
				return err
			}

			ingestStore := postgres.NewIngestRepository(db)

			before := time.Now().Add(-1 * cfg.Cron.TTL)

			return ingestStore.Delete(context.Background(), &sdump.DeleteIngestedRequestOptions{
				Before:         before,
				UseSoftDeletes: cfg.Cron.SoftDeletes,
			})
		},
	}

	rootCmd.AddCommand(cmd)
}
