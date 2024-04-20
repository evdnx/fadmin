package db

import (
	"errors"
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	DbName     string = "fadmin.db"
	AuthBucket string = "auth"
)

var db *bbolt.DB

func Init() error {
	var err error
	db, err = bbolt.Open(DbName, 0600, nil)
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

func getBucket(tx *bbolt.Tx, bucket string) (*bbolt.Bucket, error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, errors.New("bucket " + bucket + " doesn't exist")
	}

	return b, nil
}

func Read(bucket, key string) (string, error) {
	value := ""
	err := DB().View(func(tx *bbolt.Tx) error {
		b, err := getBucket(tx, bucket)
		if err != nil {
			return err
		}

		v := b.Get([]byte(key))
		if v == nil {
			return errors.New("the key " + key + " does not exist")
		}

		value = string(v)
		return nil
	})

	return value, err
}

func Update(bucket, key, value string) error {
	return DB().Update(func(tx *bbolt.Tx) error {
		b, err := getBucket(tx, bucket)
		if err != nil {
			return err
		}

		err = b.Put([]byte(key), []byte(value))
		return err
	})
}
