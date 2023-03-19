package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/comPass/pkg/hash"
	"github.com/golonzovsky/comPass/pkg/pwned"
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

			return CheckPasswordInFolder(cmd.Context(), dir, args[0])
			//return pwned.NewClient().CheckPasswordOnline(cmd.Context(), args[0])
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "~/.compass/hashes", "Output dir location")
	return cmd
}

func CheckPasswordInFolder(ctx context.Context, dir string, pass string) error {
	log.Info("Checking password on folder based storage")

	folderStore, err := storage.NewFolderStorage(dir)
	if err != nil {
		return err
	}

	hash := hash.Compute(pass)
	ensureRangeDownloaded(ctx, dir, hash[:5], folderStore)

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

func ensureRangeDownloaded(ctx context.Context, dir, prefix string, store *storage.FolderStorage) error {
	mStore, err := storage.NewMetadataStore(dir)
	if err != nil {
		return err
	}
	defer mStore.Close()

	if needsRefresh, err := mStore.NeedsRefresh(prefix); err != nil || !needsRefresh {
		return err
	}

	rangeResp, err := pwned.NewClient().DownloadRange(ctx, prefix)
	if err != nil {
		return err
	}

	err = store.SaveRange(rangeResp)
	if err != nil {
		return err
	}
	return mStore.Save(prefix, &rangeResp.RangeMetadata)
}
