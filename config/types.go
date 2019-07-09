package config

type Config struct {
	Port      string
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
