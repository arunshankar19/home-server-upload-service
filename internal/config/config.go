package config

import "github.com/kelseyhightower/envconfig"

type config struct {
	ApiVersion int `required:"true" split_words:"true"`
	ApiPort    int `required:"true" split_words:"true"`

	// db
	PostgresUser       string `required:"true" split_words:"true"`
	PostgresPassword   string `required:"true" split_words:"true"`
	PostgresHost       string `required:"true" split_words:"true"`
	PostgresPort       int    `required:"true" split_words:"true"`
	PostgresDBName     string `required:"true" split_words:"true"`
	PostgresSSLMode    string `required:"true" split_words:"true"`
	PostgresSearchPath string `required:"true" split_words:"true"`

	// observability
	ObservabilityEnabled bool   `required:"true" split_words:"true"`
	TraceExpoterURL      string `required:"true" split_words:"true"`

	// storage
	MinioAddr            string `required:"true" split_words:"true"`
	MinioAccessKeyID     string `required:"true" split_words:"true"`
	MinioSecretAccessKey string `required:"true" split_words:"true"`
	MinioDataBucketName  string `required:"true" split_words:"true"`
}

// NewAppConfig returns the app config
func NewAppConfig() (config, error) {
	cfg := config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
