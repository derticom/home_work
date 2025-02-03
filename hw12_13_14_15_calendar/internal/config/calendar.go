package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

//nolint:tagliatelle // snake_case is allowed here.
type CalendarConfig struct {
	LogLevel string `yaml:"log_level"`

	HTTPServer  HTTPServer `yaml:"http_server"`
	GRPCAddress string     `yaml:"grpc_address"`

	UseDataBaseStorage bool   `yaml:"use_data_base_storage"`
	PostgresURL        string `yaml:"postgres_url"`
}

type HTTPServer struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewCalendarConfig() *CalendarConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg CalendarConfig

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
