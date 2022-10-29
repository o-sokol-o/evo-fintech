package config

import (
	"errors"
	"os"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/o-sokol-o/evo-fintech/internal/domain"
)

func Init() *domain.Config {
	var cfg domain.Config

	setFromEnv(&cfg)

	return &cfg
}

func setFromEnv(cfg *domain.Config) error {

	if runtime.GOOS == "windows" {
		// windows specific code here...
		err := godotenv.Load()
		if err != nil {
			return errors.New("error loading .env file")
		}
	}

	cfg.Environment = os.Getenv("ENV")

	if err := envconfig.Process("db", &cfg.Postgres); err != nil {
		return err
	}
	if err := envconfig.Process("http", &cfg.HTTP); err != nil {
		return err
	}

	// cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	// cfg.HTTP.Port = os.Getenv("HTTP_PORT")

	// cfg.DB.Name = os.Getenv("DB_NAME")
	// cfg.DB.User = os.Getenv("DB_USER")
	// cfg.DB.Host = os.Getenv("DB_HOST")
	// cfg.DB.Password = os.Getenv("DB_PASSWORD")

	return nil
}
