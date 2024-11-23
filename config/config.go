package config

import (
	"github.com/kelseyhightower/envconfig"
)

func GetConfig() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

type Config struct {
	DSN   string `required:"true"`
	Token string `required:"true"`

	GraphServer Server `split_words:"true"`

	Debug bool
}

type Server struct {
	Host string ``
	Port string `default:"8080"`
}

type JWT struct {
	PrivateKey string `required:"true"`
	PublicKey  string `required:"true"`
}
