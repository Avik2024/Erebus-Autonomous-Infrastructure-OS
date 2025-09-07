package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port    string `envconfig:"PORT" default:"8080"`
	Version string `envconfig:"VERSION" default:"v0.0.1"`
	Env     string `envconfig:"ENV" default:"development"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

