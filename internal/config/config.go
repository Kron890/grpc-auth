package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// ServerPort string `yaml:"SERVER_PORT"`

	//Postgres
	PostgresPort     string `yaml:"POSTGRES_PORT"`
	PostgresUser     string `yaml:"POSTGRES_USER"`
	PostgresName     string `yaml:"POSTGRES_NAME"`
	PostgresPassword string `yaml:"POSTGRES_PASSWORD"`
	PostgresHost     string `yaml:"POSTGRES_HOST"`

	// gPRC
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     ConfigGRPC    `yaml:"grpc"`

	//Redis
	RedisPort     string `yaml:"REDIS_PORT"`
	RedisDB       int    `yaml:"REDIS_DB"`
	RedisPassword string `yaml:"REDIS_PASSWORD"`
	RedisHost     string `yaml:"REDIS_HOST"`
	RedisPoolSize int    `yaml:"REDIS_POOLSIZE"`
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
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {

	var res string

	flag.StringVar(&res, "config", "", "path to config file ")
	flag.Parse()

	if res == "" {
		panic("local.yaml not found")
		// res = os.Getenv("CONFIG_PATH")
	}
	return res
}
