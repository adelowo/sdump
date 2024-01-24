package main

import (
	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/internal/tui"
	"github.com/spf13/cobra"
)

func createSSHCommand(rootCmd *cobra.Command, cfg *config.Config) {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Start/run the TUI app",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := tui.New(cfg)

			if _, err := app.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}
