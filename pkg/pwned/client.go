package pwned

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type client struct {
	apiPrefix  string
	folderPath string
	httpClient *http.Client
}

func NewClient() *client {
	return &client{
		apiPrefix: "https://api.pwnedpasswords.com",
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *client) CheckPasswordOnline(ctx context.Context, pass string) error {
	log.Info("Checking password online")

	passSha := sha1.Sum([]byte(pass))
	passwordHash := hex.EncodeToString(passSha[:])
	hashPrefix := passwordHash[:5]

	rangeContent, err := c.DownloadRange(ctx, hashPrefix)
	if err != nil {
		return err
	}
	log.Debug("Range response:", "resp", rangeContent)

	passHashSuffix := passwordHash[5:]
	if count := occurrencesInRange(rangeContent.Body, passHashSuffix); count > 0 {
		log.Warn("Password compromised", "count", count)
		os.Exit(1)
	} else {
		log.Info("Password is safe")
	}

	return nil
}

type RangeResponse struct {
	HashPrefix         string
	Body               []byte
	ETag               string
	CloudflareCacheHit bool
	LastModified       time.Time
	// "?ntlm=true"
}

func (rr *RangeResponse) String() string {
	return fmt.Sprintf("HashPrefix: %s, BodyLen: %d, ETag: %s, CloudflareCacheHit: %t, LastModified: %s",
		rr.HashPrefix, len(rr.Body), rr.ETag, rr.CloudflareCacheHit, rr.LastModified)
}

// todo use If-None-Match header to check if changed in case its overide, take etag from metadata store
func (c *client) DownloadRange(ctx context.Context, hashPrefix string) (*RangeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.apiPrefix+"/range/"+hashPrefix, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lm := resp.Header.Get("Last-Modified")
	lastModified, err := time.Parse(time.RFC1123, lm)

	return &RangeResponse{
		HashPrefix:         hashPrefix,
		Body:               responseData,
		ETag:               resp.Header.Get("ETag"),
		CloudflareCacheHit: resp.Header.Get("CF-Cache-Status") == "HIT",
		LastModified:       lastModified,
	}, nil
}

func occurrencesInRange(rangeResp []byte, passHashSuffix string) int {
	for _, line := range strings.Split(string(rangeResp), "\r\n") {
		split := strings.Split(line, ":")
		hashSuffix := strings.ToLower(split[0])
		count := split[1]
		if hashSuffix == passHashSuffix {
			occurrences, _ := strconv.Atoi(count)
			return occurrences
		}
	}
	return 0
}