package auth

import (
	"fileserver/config"
	"fileserver/db"
	"fileserver/types"
	"fileserver/user"
	"strings"

	"github.com/casbin/casbin"
	boltadapter "github.com/wirepair/bolt-adapter"
)

type auth struct {
	*casbin.Enforcer
}

var instance *auth

func New(modelPath string) {
	adapter := boltadapter.NewAdapter(db.GetInstance().DB)

	cas := casbin.NewEnforcer(modelPath, adapter)
	instance = &auth{
		cas,
	}
	instance.AddFunction("prefixMatch", PrefixMatchFunc)

	instance.LoadPolicy()
	instance.initData()
}

func GetInstance() *auth {
	return instance
}

func (a *auth) SaveAndReloadPolicy() (err error) {
	err = a.SavePolicy()
	if err != nil {
		return err
	}
	err = a.LoadPolicy()
	if err != nil {
		return err
	}

	return
}

func (a *auth) initData() {
	cfg := config.GetInstance()
	userM := user.GetInstance()

	roleAdmin := "role_admin"
	roleDefault := "role_default"

	admin := &types.User{
		Name:     cfg.Admin.Username,
		Password: cfg.Admin.Password,
	}
	user, err := userM.Login(admin)
	if err != nil {
		_ = userM.Register(admin)
		user, err = userM.Login(admin)
	}
	policy := &types.Policy{
		Role:   roleAdmin,
		Path:   `/*`,
		Method: ".*",
	}

	a.AddPermissionForUser(policy.Role, policy.Path, policy.Method)
	a.AddRoleForUser(user.ID.String(), roleAdmin)
	a.AddRoleForUser("0", roleDefault)
	a.SaveAndReloadPolicy()
}

func PrefixMatch(key1 string, key2 string) bool {
	i := strings.Index(key2, "*")
	if i == -1 {
		return strings.HasPrefix(key2, key1)
	}

	if len(key1) > i {
		return strings.HasPrefix(key2[:i], key1[:i])
	}
	return strings.HasPrefix(key2[:i], key1)
}

func PrefixMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(PrefixMatch(name1, name2)), nil
}
