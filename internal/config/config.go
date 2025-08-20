package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ServerPort string        `yaml:"SERVER_PORT"`
	DBPort     string        `yaml:"POSTGRES_PORT"`
	User       string        `yaml:"POSTGRES_USER"`
	DB         string        `yaml:"POSTGRES_DB"`
	DBPassword string        `yaml:"POSTGRES_DBPASSWORD"`
	Host       string        `yaml:"POSTGRES_HOST"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC       ConfigGRPC    `yaml:"grpc"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Загрузка конфига
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {

	var res string

	flag.StringVar(&res, "config", "", "path to config file ")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
