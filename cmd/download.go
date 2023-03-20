package cmd

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/compass/pkg/hash"
	"github.com/golonzovsky/compass/pkg/pwned"
	"github.com/golonzovsky/compass/pkg/storage"
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

			mdStore, err := storage.NewMetadataStore(options.OutDir)
			if err != nil {
				return err
			}
			defer mdStore.Close()

			store, err := storage.NewFolderStorage(options.OutDir)
			if err != nil {
				return err
			}

			// todo use len(channel) to find out which side is saturated in order to find out optimal parallelism

			hashCh := hash.NewPrefixGen()
			hashes, _ := downloadHashes(cmd.Context(), mdStore, hashCh, options.Parallel)
			done, _ := storeRanges(cmd.Context(), store, mdStore, hashes, options.Parallel)

			<-done
			// print stats
			// "Finished downloading all hash ranges in ElapsedMilliseconds-ms (HashesPerSecond hashes per second)."
			// "We made CloudflareRequests Cloudflare requests (avg response time: CloudflareRequestTimeTotal/CloudflareRequests-ms). Of those, Cloudflare had already cached CloudflareHits requests, and made CloudflareMisses requests to the Have I Been Pwned origin server."

			return nil
		},
	}
	cmd.Flags().StringVar(&options.OutDir, "dir", "~/.compass/hashes", "Output dir location")
	cmd.Flags().IntVarP(&options.Parallel, "parallel", "p", 100, "Number of parallel downloads")

	return cmd
}

func storeRanges(ctx context.Context, store *storage.FolderStorage, mdStore *storage.MetadataStore, rangeRespCh <-chan *pwned.RangeResponse, parallel int) (<-chan struct{}, <-chan error) {
	errs := make(chan error, 1)
	done := make(chan struct{})

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < parallel; i++ {
		g.Go(func() error {
			for rr := range rangeRespCh {
				if err := store.SaveRange(rr); err != nil {
					return err
				}
				if err := mdStore.Save(rr.HashPrefix, &rr.RangeMetadata); err != nil {
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

func downloadHashes(ctx context.Context, metadataStore *storage.MetadataStore, hashCh <-chan string, parallel int) (<-chan *pwned.RangeResponse, <-chan error) {
	ranges := make(chan *pwned.RangeResponse, parallel)
	errs := make(chan error, 1)
	pwndClient := pwned.NewClient()
	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < parallel; i++ {
		g.Go(func() error {
			for prefix := range hashCh {
				refresh, err := metadataStore.NeedsRefresh(prefix)
				if err != nil {
					return err
				}
				if !refresh {
					log.Debug("Skipping hash range, as its to date", "hashPrefix", prefix)
					continue
				}
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
