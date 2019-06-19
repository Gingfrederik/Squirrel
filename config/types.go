package config

type Config struct {
	Root      string
	DB        string
	SecretKey string
	Casbin    struct {
		Model string
	}
	Admin struct {
		Username string
		Password string
	}
}
