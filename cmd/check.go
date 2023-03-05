package cmd

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

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
			// todo check online

			passSha := sha1.Sum([]byte(args[0]))
			passwordHash := hex.EncodeToString(passSha[:])
			resp, err := http.Get("https://api.pwnedpasswords.com/range/" + passwordHash[:5])
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
				hash := strings.ToLower(split[0])
				count := split[1]
				if hash == passHashSuffix {
					log.Warn("Password compromised", "count", count)
					return nil
				}
			}
			log.Info("Password is safe")
			return nil
		},
	}

	return cmd
}
