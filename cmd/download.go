package cmd

import (
	"github.com/golonzovsky/compass/pkg/download"
	"github.com/spf13/cobra"
)

func NewDownloadCmd() *cobra.Command {
	var options struct {
		OutDir   string
		Parallel int
	}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "download hashes from haveibeenpwned.com (~25Gb)",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, closeStore, err := download.NewDownloader(options.OutDir, options.Parallel)
			if err != nil {
				return err
			}
			defer closeStore()
			return d.Download(cmd.Context())
		},
	}
	cmd.Flags().StringVar(&options.OutDir, "dir", "~/.compass/hashes", "Output dir location")
	cmd.Flags().IntVarP(&options.Parallel, "parallel", "p", 100, "Number of parallel downloads")

	return cmd
}
