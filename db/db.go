package db

import (
	"fmt"

	"github.com/evdnx/unixmint/constants"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB

func Init() error {
	var err error
	db, err = bbolt.Open(constants.DbName, 0600, nil)
	if err != nil {
		return err
	}

	// create initial structure
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte("auth"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	return nil
}

func DB() *bbolt.DB {
	return db
}
