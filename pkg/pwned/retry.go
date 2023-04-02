package pwned

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

// todo add tests
func (c *client) doWithRetry(ctx context.Context, req *http.Request, retries int) (*http.Response, error) {
	req = req.WithContext(ctx)
	base, cap := 500*time.Millisecond, 30*time.Second
	resp, err := c.httpClient.Do(req)
	r := 0

	// todo do check err and filter out unrecoverable ones
	for backoff := base; err != nil && r <= retries; backoff <<= 1 {
		r++
		if backoff > cap {
			backoff = cap
		}

		jitter := rand.Int63n(int64(backoff * 3))
		sleep := base + time.Duration(jitter)

		log.Debug("Request attempt failed", "attempt", r, "delay", sleep, "err", err)

		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		resp, err = c.httpClient.Do(req)
	}
	return resp, err
}
