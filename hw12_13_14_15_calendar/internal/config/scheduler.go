package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

//nolint:tagliatelle // snake_case is allowed here.
type SchedulerConfig struct {
	LogLevel string `yaml:"log_level"`

	AmpqURI     string `yaml:"amqp_uri"`
	GRPCAddress string `yaml:"grpc_address"`

	UseDataBaseStorage bool          `yaml:"use_data_base_storage"`
	ProcessPeriod      time.Duration `yaml:"process_period"`
}

func NewSchedulerConfig() *SchedulerConfig {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg SchedulerConfig

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
