package db

import (
	"os"
	"sync"

	"github.com/boltdb/bolt"
)

const (
	userBucket     string = "user"
	userNameBucket string = "user.name"
)

type db struct {
	*bolt.DB
}

var instance *db
var once sync.Once

func New(path string) {
	once.Do(func() {
		instance = &db{}

		if err := instance.Open(path, 0600); err != nil {
			panic(err)
		}

	})
}

func GetInstance() *db {
	return instance
}

// Open initializes and opens the database.
func (db *db) Open(path string, mode os.FileMode) error {
	var err error

	db.DB, err = bolt.Open(path, mode, nil)
	if err != nil {
		return err
	}

	// Create User buckets.
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(userBucket)); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		db.Close()
		return err
	}

	// Create Name ID Buckets
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(userNameBucket)); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		db.Close()
		return err
	}

	return nil
}
