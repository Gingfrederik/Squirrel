package config

import (
	"log"

	"github.com/spf13/viper"
)

var configuration *Config

func New() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetDefault("Root", ".")
	viper.SetDefault("Admin.Username", "root")
	viper.SetDefault("Admin.Password", "1234")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode config into struct, %v", err)
	}

	return configuration
}

func GetInstance() *Config {
	return configuration
}
