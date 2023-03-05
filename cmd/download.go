package cmd

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/comPass/pkg/hash"
	"github.com/golonzovsky/comPass/pkg/pwned"
	"github.com/golonzovsky/comPass/pkg/storage"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type DownloadStats struct {
	CloudflareRequests          int
	CloudflareHits              int
	CloudflareMisses            int
	CloudflareRequestsTimeTotal time.Duration
	CloudFlareHitPercentage     float64
	CloudFlareMissPercentage    float64
	CloudflareRequestsPerSecond float64
	HashesDownloaded            int
	HashesPerSecond             float64
	ElapsedTime                 time.Duration
}

func NewDownloadCmd() *cobra.Command {
	var options struct {
		OutDir   string
		Parallel int
	}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "download hashes from haveibeenpwned.com (~25Gb)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo update stats every 100ms

			store, err := storage.NewFolderStorage(options.OutDir)
			if err != nil {
				return err
			}

			// todo use len(channel) to find out which side is saturated in order to find out optimal parallelism

			hashCh := hash.NewPrefixGen()
			hashes, _ := downloadHashes(cmd.Context(), hashCh, options.Parallel)
			done, _ := storeRanges(cmd.Context(), store, hashes, options.Parallel)

			<-done
			// print stats
			// "Finished downloading all hash ranges in ElapsedMilliseconds-ms (HashesPerSecond hashes per second)."
			// "We made CloudflareRequests Cloudflare requests (avg response time: CloudflareRequestTimeTotal/CloudflareRequests-ms). Of those, Cloudflare had already cached CloudflareHits requests, and made CloudflareMisses requests to the Have I Been Pwned origin server."

			return nil
		},
	}
	cmd.Flags().StringVarP(&options.OutDir, "output-dir", "o", "~/.compass/hashes", "Output dir location")
	cmd.Flags().IntVarP(&options.Parallel, "parallel", "p", 10, "Number of parallel downloads")

	return cmd
}

func storeRanges(ctx context.Context, store *storage.FolderStorage, rangeRespCh <-chan *pwned.RangeResponse, parallel int) (<-chan struct{}, <-chan error) {
	errs := make(chan error, 1)
	done := make(chan struct{})

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < parallel; i++ {
		g.Go(func() error {
			for rr := range rangeRespCh {
				err := store.SaveRange(ctx, rr)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	go func() {
		err := g.Wait()
		if err != nil && err != context.Canceled {
			log.Error("Failure during hash store", "err", err.Error())
			errs <- err
			close(errs)
		}
		close(done)
	}()

	return done, errs
}

func downloadHashes(ctx context.Context, hashCh <-chan string, parallel int) (<-chan *pwned.RangeResponse, <-chan error) {
	ranges := make(chan *pwned.RangeResponse, parallel)
	errs := make(chan error, 1)
	pwndClient := pwned.NewClient()
	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < parallel; i++ {
		g.Go(func() error {
			for prefix := range hashCh {
				rangeResp, err := pwndClient.DownloadRange(ctx, prefix)
				if err != nil {
					return err
				}
				ranges <- rangeResp
			}
			return nil
		})
	}
	go func() {
		err := g.Wait()
		if err != nil && err != context.Canceled {
			log.Error("Failure during hash download", "err", err.Error())
			errs <- err
			close(errs)
		}
		close(ranges)
	}()
	return ranges, errs
}
