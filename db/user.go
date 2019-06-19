package db

import (
	"encoding/json"
	"errors"

	"github.com/bwmarrin/snowflake"

	"github.com/boltdb/bolt"

	"fileserver/types"
)

var (
	ErrKeyNotExists = errors.New("key not exists")
)

func (db *db) AddOrUpdateUser(user *types.User) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = b.Put(user.ID.Bytes(), buf)
		if err != nil {
			return err
		}

		b = tx.Bucket([]byte(userNameBucket))

		err = b.Put([]byte(user.Name), user.ID.Bytes())
		if err != nil {
			return err
		}

		return nil
	})
}

func (db *db) GetUserByID(id snowflake.ID) (user *types.User, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		v := b.Get(id.Bytes())
		if v == nil {
			return ErrKeyNotExists
		}

		user, err = decodeToUser(v)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (db *db) GetAllUser() (user []*types.UserNoPass, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))

		b.ForEach(func(k, v []byte) error {
			u, err := decodeToUserNoPass(v)
			if err != nil {
				return err
			}

			user = append(user, u)
			return nil
		})
		return nil
	})

	return
}

func (db *db) GetIDByName(name string) (id snowflake.ID, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userNameBucket))
		v := b.Get([]byte(name))
		if v == nil {
			return ErrKeyNotExists
		}

		id, err = snowflake.ParseBytes(v)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func decodeToUser(data []byte) (user *types.User, err error) {
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	return
}

func decodeToUserNoPass(data []byte) (user *types.UserNoPass, err error) {
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	return
}
