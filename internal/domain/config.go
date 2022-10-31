package domain

import (
	"github.com/o-sokol-o/evo-fintech/pkg/database"
)

type HTTPConfig struct {
	Host string
	Port string
}

type Config struct {
	Postgres database.PostgresConfig
	HTTP     HTTPConfig
}
