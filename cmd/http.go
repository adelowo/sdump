package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/server/httpd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func createHTTPCommand(cmd *cobra.Command, cfg *config.Config) {
	command := &cobra.Command{
		Use:  "http",
		Long: "Start/run the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			hostName, err := os.Hostname()
			if err != nil {
				return err
			}

			logger := logrus.WithField("host", hostName).
				WithField("module", "http.server")

			httpServer := httpd.New(*cfg, nil, logger)

			go func() {
				if err := httpServer.ListenAndServe(); err != nil {
					logger.WithError(err).Fatal("could not start http server")
				}
			}()

			<-sig

			return nil
		},
	}

	cmd.AddCommand(command)
}
