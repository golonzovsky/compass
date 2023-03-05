package cmd

import (
	"github.com/golonzovsky/comPass/pkg/pwned"
	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "check password. Uses Bloom filters, local hashes when available, otherwise checks online",
		ValidArgsFunction: cobra.NoFileCompletions,
		Args:              cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// todo check if bloom available and evaluate

			// todo check if local file based hashes available and evaluate

			// todo check if local dir based hashes available and evaluate

			return pwned.NewClient().CheckPasswordOnline(cmd.Context(), args[0])
		},
	}

	return cmd
}
