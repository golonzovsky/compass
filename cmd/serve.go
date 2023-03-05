package cmd

import (
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve rest api",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo implement
			return nil
		},
	}

	return cmd
}
