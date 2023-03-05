package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/comPass/pkg/pwned"
)

// todo implement metadata (etag, lastModified etc) storage as well
// todo sharded store

type FolderStorage struct {
	path string
}

func NewFolderStorage(path string) (*FolderStorage, error) {
	if strings.HasPrefix(path, "~/") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = homedir + path[1:]
	}

	err := ensureFolderExists(path)
	if err != nil {
		return nil, err
	}
	return &FolderStorage{path: path}, nil
}

func ensureFolderExists(path string) error {
	dir, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	if err != nil {
		return err
	}
	if !dir.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	//_, err = f.Readdirnames(1) // Or f.Readdir(1)
	//if err != io.EOF {
	//	return fmt.Errorf("folder %s is not empty", path)
	//}

	return nil
}

func (fs *FolderStorage) SaveRange(ctx context.Context, rr *pwned.RangeResponse) error {
	log.Debug("Saving range", "hashPrefix", rr.HashPrefix)
	f, err := os.Create(fs.path + "/" + rr.HashPrefix + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(rr.Body)
	return err
}
