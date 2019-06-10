package config

type Config struct {
	Root  string
	Admin struct {
		Username string
		Password string
	}
}
