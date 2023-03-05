package cmd

import (
	"time"

	"github.com/spf13/cobra"
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
		output     string
		parallel   int
		override   bool
		singleFile bool
		ntlm       bool
	}

	cmd := &cobra.Command{
		Use:   "download",
		Short: "download hashes from haveibeenpwned.com (~25Gb)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo check if output file exists and fail if not --override or delete
			// todo check if output dir exists and not empty and fail if not --override or delete
			// todo update stats every 100ms

			// todo do download and store location and etags in ~/.compass/last-download.yml

			// client
			// todo retry Async 10 OnRequestError
			// tls 1.3, http2
			// https://api.pwnedpasswords.com/range/
			// UserAgent: hibp-go-downloader + version

			// print stats
			// "Finished downloading all hash ranges in ElapsedMilliseconds-ms (HashesPerSecond hashes per second)."
			// "We made CloudflareRequests Cloudflare requests (avg response time: CloudflareRequestTimeTotal/CloudflareRequests-ms). Of those, Cloudflare had already cached CloudflareHits requests, and made CloudflareMisses requests to the Have I Been Pwned origin server."

			return nil
		},
	}
	cmd.Flags().StringVarP(&options.output, "output", "o", "~/.compass/hashes.txt", "Output file")
	cmd.Flags().BoolVarP(&options.singleFile, "single-file", "s", true, "Download all hashes to a single file")
	cmd.Flags().BoolVarP(&options.override, "override", "O", false, "Override output if exists")
	cmd.Flags().IntVarP(&options.parallel, "parallel", "p", 10, "Number of parallel downloads")
	cmd.Flags().BoolVarP(&options.ntlm, "ntlm", "n", false, "Download NTLM hashes instead of SHA1")

	return cmd
}

func getHashRange(i int, ntlm bool) []byte {
	// todo start stopwatch
	requestUri := getHashRangeUrl(i)
	if ntlm {
		requestUri += "?ntlm=true"
	}
	// todo make request
	// use If-None-Match header to check if changed in case its overide

	// todo stop stopwatch
	// todo update stats

	// resp.headers["CF-Cache-Status"] == "HIT" => CloudflareHits++ else CloudflareMisses++

	return nil
}

func getHashRangeUrl(i int) string {
	// todo get hash range url
	return ""
}