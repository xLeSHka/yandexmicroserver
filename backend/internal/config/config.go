package config

import (
	"log"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// internal/config/config.go

type Config struct {
	HTTPServer    `yaml:"http_server"`
	StorageConfig `yaml:"storage"`
}
type StorageConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	MaxAttempts int    `json:"maxattempts"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0:8080"`
}

var instance *Config
var once sync.Once

// internal/config/config.go

func GetConfig(logger slog.Logger) *Config {
	once.Do(func() {
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig(".\\backend\\Config\\config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			log.Fatalf("Configuration error: %s", err)
		}
	})
	return instance
}
