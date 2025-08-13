package config

import (
	"auth-grpc/internal/envcfg"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBPort     string `env:"POSTGRES_PORT,required"`
	User       string `env:"POSTGRES_USER,required"`
	DB         string `env:"POSTGRES_DB,required"`
	Host       string `env:"POSTGRES_HOST,required"`
}

type ConfigGRPC struct {
	Port    string `env:"port"`
	Timeout string `env:"timeout"`
}

func GetConfigGRPC() (ConfigGRPC, error) {
	if err := godotenv.Load(); err != nil {
		return ConfigGRPC{}, err
	}
	var cfg ConfigGRPC

	cfg.Port = os.Getenv("GRPC_PORT")
	if cfg.Port == "" {
		return ConfigGRPC{}, errors.New("GRPC_PORT not found")
	}
	cfg.Timeout = os.Getenv("GRPC_TIMEOUT")
	if cfg.Timeout == "" {
		return ConfigGRPC{}, errors.New("GRPC_TIMEOUT not found")
	}
	return cfg, nil

}

func GetConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := envcfg.Load(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
