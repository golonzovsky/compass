package cmd

import (
	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "check",
		Short: "check password. Uses Bloom filters, local hashes when available, otherwise checks online",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return cmd
}
