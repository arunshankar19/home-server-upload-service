package config

import "github.com/kelseyhightower/envconfig"

type config struct {
	ApiVersion int `required:"true" split_words:"true"`
	ApiPort    int `required:"true" split_words:"true"`
}

func NewAppConfig() (config, error) {
	cfg := config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
