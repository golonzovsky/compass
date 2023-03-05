package cmd

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "check",
		Short:             "check password. Uses Bloom filters, local hashes when available, otherwise checks online",
		ValidArgsFunction: cobra.NoFileCompletions,
		Args:              cobra.MatchAll(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo check if bloom available and evaluate

			// todo check if local hashes available and evaluate
			cmd.SilenceUsage = true
			return checkPasswordOnline(cmd.Context(), args[0])
		},
	}

	return cmd
}

func checkPasswordOnline(ctx context.Context, pass string) error {
	log.Info("Checking password online")

	passSha := sha1.Sum([]byte(pass))
	passwordHash := hex.EncodeToString(passSha[:])
	hashPrefix := passwordHash[:5]

	reqContext, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqContext, "GET", "https://api.pwnedpasswords.com/range/"+hashPrefix, nil)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	passHashSuffix := passwordHash[5:]
	for _, line := range strings.Split(string(responseData), "\r\n") {
		split := strings.Split(line, ":")
		hashSuffix := strings.ToLower(split[0])
		count := split[1]
		if hashSuffix == passHashSuffix {
			log.Warn("Password compromised", "count", count)
			return nil
		}
	}
	log.Info("Password is safe")
	return nil
}
