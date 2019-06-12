package config

type Config struct {
	Root      string
	DB        string
	SecretKey string
	Admin     struct {
		Username string
		Password string
	}
}
