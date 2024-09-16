package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"time"
)

type Config struct {
	Address     string        `env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	PostgresConfig
}

type PostgresConfig struct {
	ConnURL      string `env:"POSTGRES_CONN"`
	JDBCURL      string `env:"POSTGRES_JDBC_URL"`
	Username     string `env:"POSTGRES_USERNAME`
	Password     string `env:"POSTGRES_PASSWORD"`
	Host         string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port         int    `env:"POSTGRES_PORT" env-default:"5432"`
	DatabaseName string `env:"POSTGRES_DATABASE"`
}

func MustLoad() *Config {
	var cfg Config

	if os.Getenv("ENV") == "local" {
		log.Println("local")
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("cannot pars vars to a config struct", err)
	}

	log.Println(cfg)
	return &cfg
}
