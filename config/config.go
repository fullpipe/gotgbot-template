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

	Debug bool
}

type JWT struct {
	PrivateKey string `required:"true"`
	PublicKey  string `required:"true"`
}
