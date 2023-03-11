package cmd

import (
	"github.com/spf13/cobra"
)

func NewBloomCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "bloom",
		Short: "build bloom filters from password hashes",
		RunE: func(cmd *cobra.Command, args []string) error {
			//todo implement
			return nil
		},
	}

	return cmd
}
