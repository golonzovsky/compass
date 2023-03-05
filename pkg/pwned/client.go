package pwned

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type client struct {
	apiPrefix  string
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

	req, err := http.NewRequestWithContext(ctx, "GET", c.apiPrefix+"/range/"+hashPrefix, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	passHashSuffix := passwordHash[5:]
	if count := containsInRange(responseData, passHashSuffix); count > 0 {
		log.Warn("Password compromised", "count", count)
	} else {
		log.Info("Password is safe")
	}

	return nil
}

func containsInRange(rangeResp []byte, passHashSuffix string) int {
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
