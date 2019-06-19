package user

import (
	"errors"
	"fileserver/db"
	"fileserver/types"
	"time"

	"github.com/bwmarrin/snowflake"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	genID *snowflake.Node
}

var instance *user

func New() {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	instance = &user{
		genID: node,
	}
}

func GetInstance() *user {
	return instance
}

func (u *user) Register(user *types.User) (err error) {
	dbAgent := db.GetInstance()
	_, err = dbAgent.GetIDByName(user.Name)
	if err != db.ErrKeyNotExists {
		err = errors.New("user exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user.ID = u.genID.Generate()
	user.Password = string(hash)
	user.Locked = false
	user.CreatedAt = time.Now().UTC()
	user.UpdateAt = time.Now().UTC()

	err = dbAgent.AddOrUpdateUser(user)
	if err != nil {
		return
	}

	return
}

func (u *user) Login(user *types.User) (dbUser *types.User, err error) {
	dbAgent := db.GetInstance()
	id, err := dbAgent.GetIDByName(user.Name)
	if err != nil {
		return
	}

	dbUser, err = dbAgent.GetUserByID(id)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return
	}

	return
}
