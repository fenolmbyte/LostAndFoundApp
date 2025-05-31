package server_config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	CONFIG_SERVER_PATH = "CONFIG_SERVER_PATH"
)

type Config struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadServerConfig() (*Config, error) {

	slog.Debug("Loading server config")

	configPath := os.Getenv(CONFIG_SERVER_PATH)
	if configPath == "" {
		return nil, fmt.Errorf("%s environment variable not set", CONFIG_SERVER_PATH)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist %s", CONFIG_SERVER_PATH, configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("cannot load config file: %s", err)
	}

	return &config, nil
}
