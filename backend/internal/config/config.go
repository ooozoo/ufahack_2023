package config

import "github.com/ilyakaznacheev/cleanenv"

type ServerConfig struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"localhost"`
	Port    string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env-default:"user"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"secretpass"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Description() (string, error) {
	return cleanenv.GetDescription(&Config{}, nil)
}
