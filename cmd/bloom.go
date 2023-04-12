package cmd

import (
	"encoding/json"
	"os"

	"github.com/golonzovsky/compass/pkg/bloom"
	"github.com/golonzovsky/compass/pkg/storage"
	"github.com/spf13/cobra"
)

func NewBloomCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "bloom",
		Short: "build bloom filters from password hashes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doBuildBloomFilters(dir)
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "~/.compass/hashes", "Output dir location")

	return cmd
}

func doBuildBloomFilters(dir string) error {
	ms, closeMd, err := storage.NewMetadataStore(dir)
	if err != nil {
		return err
	}
	defer closeMd()

	keys, err := ms.StreamAvailableKeys()
	if err != nil {
		return err
	}

	fs, err := storage.NewFolderStorage(dir)
	if err != nil {
		return err
	}

	bf := bloom.NewFilter(600000000, 0.000001)
	for k := range keys {
		hashes, err := fs.GetRange(k)
		if err != nil {
			return err
		}
		for _, h := range hashes {
			bf.Add([]byte(h))
		}
	}

	bfMarshal, err := json.Marshal(bf)
	if err != nil {
		return err
	}

	f, err := os.Create(dir + "/" + "bloom_filter.json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bfMarshal)

	return nil
}
