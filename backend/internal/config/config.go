package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Address     string        `yaml:"address" env:"SERVER_ADDRESS" env-default:"localhost"`
	Port        int           `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-default:"30m"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"SERVER_IDLE_TIMEOUT" env-default:"30m"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env-default:"user"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret" env:"JWT_SECRET"`
	TTL    time.Duration `yaml:"ttl" env:"JWT_TTL" env-default:"30m"`
}

type Config struct {
	Env      string         `yaml:"env" env:"ENV" env-default:"local"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	f := flag.NewFlagSet("Configuration", flag.ExitOnError)
	f.StringVar(&res, "config", "", "Path to config file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(&Config{}, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
