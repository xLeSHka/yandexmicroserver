package config

import (
	"log"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

// internal/config/config.go

type Config struct {
	HTTPServer    `yaml:"httpmicroserver"`
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
	Host       string `yaml:"host" env-default:"0.0.0.0"`
	OrchPort   string `yaml:"orch_port" env-default:"8082"`
	AgentPort  string `yaml:"agent_port" env-default:"3030"`
	AgentPort1 string `yaml:"agent_port1" env-default:"3031"`
	AgentPort2 string `yaml:"agent_port2" env-default:"3032"`
	AgentPort3 string `yaml:"agent_port3" env-default:"3033"`
	AgentPort4 string `yaml:"agent_port4" env-default:"3034"`
	AgentCount string `yaml:"agent_count" env-default:"5"`
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
