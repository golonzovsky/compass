package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestSandbox(t *testing.T) {
	err := doTest()
	require.NoError(t, err)
}

func doTest() error {
	db, err := bolt.Open("test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		// readwrite transaction
		_, err := tx.CreateBucketIfNotExists([]byte("RangeETag"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("RangeETag"))
		return b.Put([]byte("answer"), []byte("42"))
	})
	if err != nil {
		return err
	}

	err = db.View(func(tx *bolt.Tx) error {
		// readonly transaction
		return nil
	})
	if err != nil {
		return err
	}

	return err
}

func manualTransaction(db *bolt.DB) error {
	// Start a writable transaction.
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err = tx.CreateBucket([]byte("MyBucket"))
	if err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
