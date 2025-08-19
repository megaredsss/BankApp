package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string `mapstructure:"env" env-default:"dev" env-required:"true"`
	HTTP_Server `mapstructure:"http_server"`
	Database    `mapstructure:"database"`
	SecretKey   []byte `mapstructure:"secret_key" env-required:"true"`
}

type HTTP_Server struct {
	Host         string `mapstructure:"host" env-default:"localhost:8080"`
	Timeout      int    `mapstructure:"timeout"`
	Idle_timeout int    `mapstructure:"idle_timeout"`
}
type Database struct {
	Host     string `mapstructuretructure:"host" env-default:"localhost"`
	Port     int    `mapstructure:"port" env-default:"5432"`
	User     string `mapstructure:"user" env-required:"true"`
	Password string `mapstructure:"pass" env-required:"true"`
	Name     string `mapstructure:"name" env-required:"true"`
}

func Loader() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path is not set in environment variables")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Conf&ig file does not exist at path: %s", configPath)
	}

	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	return &cfg
}
