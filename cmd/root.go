package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compass",
		Short: "compass is a tool to download/check compromised password hashes from haveibeenpwned.com",
	}
	cmd.AddCommand(NewDownloadCmd())
	cmd.AddCommand(NewBloomCmd())
	cmd.AddCommand(NewCheckCmd())
	return cmd
}
