package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golonzovsky/comPass/pkg/hash"
	"github.com/golonzovsky/comPass/pkg/pwned"
	bolt "go.etcd.io/bbolt"
)

type MetadataStore struct {
	db *bolt.DB
}

func NewMetadataStore(path string) (*MetadataStore, error) {
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

	db, err := bolt.Open(path+"/metadata.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(rangeMetadataBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &MetadataStore{
		db: db,
	}, nil
}

func (ms *MetadataStore) Save(hashPrefix string, metadata *pwned.RangeMetadata) error {
	m, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	return ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(rangeMetadataBucket)
		hashByte, err := hash.FromPrefix(hashPrefix)
		if err != nil {
			return err
		}
		return b.Put(hashByte, m)
	})
}

func (ms *MetadataStore) get(hashPrefix string) (*pwned.RangeMetadata, error) {
	var metadata pwned.RangeMetadata
	return &metadata, ms.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(rangeMetadataBucket)
		hashByte, err := hash.FromPrefix(hashPrefix)
		if err != nil {
			return err
		}
		m := b.Get(hashByte)
		if m == nil {
			return nil
		}
		return json.Unmarshal(m, &metadata)
	})
}

func (ms *MetadataStore) NeedsRefresh(hashPrefix string) (bool, error) {
	metadata, err := ms.get(hashPrefix)
	if err != nil {
		return false, err
	}
	if metadata == nil {
		return true, nil
	}
	return time.Now().After(metadata.Expires), nil
}

func (ms *MetadataStore) Close() {
	ms.db.Close()
}
