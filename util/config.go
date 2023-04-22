package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment  string `mapstructure:"ENVIRONMENT"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBName       string `mapstructure:"DB_NAME"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`

	DBDockerPort   string `mapstructure:"DB_DOCKER_PORT"`
	HTTPServerPort string `mapstructure:"GO_DOCKER_PORT"`

	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	AllowedServer string `mapstructure:"ALLOWED_SERVER"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
