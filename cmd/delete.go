package main

import (
	"github.com/spf13/cobra"
)

func createDeleteCommand(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "delete-http",
		Aliases: []string{"d"},
		Short:   "Deletes all old HTTP requests to preserve DB space",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	rootCmd.AddCommand(cmd)
}
