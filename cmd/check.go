package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/comPass/pkg/storage"
	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "check password. Uses Bloom filters, local hashes when available, otherwise checks online",
		ValidArgsFunction: cobra.NoFileCompletions,
		Args:              cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			// todo check if bloom available and evaluate

			// todo check if local file based hashes available and evaluate

			return CheckPasswordInFolder(dir, args[0])
			//return pwned.NewClient().CheckPasswordOnline(cmd.Context(), args[0])
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "~/.compass/hashes", "Output dir location")
	return cmd
}

func CheckPasswordInFolder(dir string, pass string) error {
	log.Info("Checking password on folder based storage")

	folderStore, err := storage.NewFolderStorage(dir)
	if err != nil {
		return err
	}
	hash := pass // todo compute hash on input hash.Compute(pass)
	contains, err := folderStore.Contains(hash)
	if err != nil {
		return err
	}
	if contains {
		log.Warn("Password compromised")
		os.Exit(1)
	} else {
		log.Info("Password is safe")
	}
	return nil
}
