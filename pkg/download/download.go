package download

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/compass/pkg/hash"
	"github.com/golonzovsky/compass/pkg/pwned"
	"github.com/golonzovsky/compass/pkg/storage"
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

type downloader struct {
	mdStore *storage.MetadataStore
	store   *storage.FolderStorage

	parallel   int
	pwndClient pwned.Client
}

func NewDownloader(outDir string, parallel int) (*downloader, func(), error) {
	noop := func() {}
	mdStore, err := storage.NewMetadataStore(outDir)
	if err != nil {
		return nil, noop, err
	}

	store, err := storage.NewFolderStorage(outDir)
	if err != nil {
		return nil, noop, err
	}

	pwndClient := pwned.NewClient()

	return &downloader{
			mdStore:    mdStore,
			store:      store,
			pwndClient: pwndClient,
			parallel:   parallel,
		}, func() {
			mdStore.Close()
		}, nil
}

func (d *downloader) Download(ctx context.Context) error {
	hashCh := hash.NewPrefixGen()
	hashes, _ := d.downloadHashes(ctx, hashCh)
	done, _ := d.storeRanges(ctx, hashes)

	<-done
	return nil
}

func (d *downloader) downloadHashes(ctx context.Context, hashCh <-chan string) (<-chan *pwned.RangeResponse, <-chan error) {
	ranges := make(chan *pwned.RangeResponse, d.parallel)
	errs := make(chan error, 1)

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < d.parallel; i++ {
		g.Go(func() error {
			for prefix := range hashCh {
				refresh, err := d.mdStore.NeedsRefresh(prefix)
				if err != nil {
					return err
				}
				if !refresh {
					log.Debug("Skipping hash range, as its up to date", "hashPrefix", prefix)
					continue
				}
				rangeResp, err := d.pwndClient.DownloadRange(ctx, prefix)
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

func (d *downloader) storeRanges(ctx context.Context, rangeRespCh <-chan *pwned.RangeResponse) (<-chan struct{}, <-chan error) {
	errs := make(chan error, 1)
	done := make(chan struct{})

	g, _ := errgroup.WithContext(ctx)
	for i := 0; i < d.parallel; i++ {
		g.Go(func() error {
			for rr := range rangeRespCh {
				if err := d.store.SaveRange(rr); err != nil {
					return err
				}
				if err := d.mdStore.Save(rr.HashPrefix, &rr.RangeMetadata); err != nil {
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
