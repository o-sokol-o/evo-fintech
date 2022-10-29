package domain

import (
	"time"

	"github.com/o-sokol-o/evo-fintech/pkg/database"
)

type HTTPConfig struct {
	Host string
	Port string
}

type NetServerConfig struct {
	Host string
	Port int
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Salt            string
	Secret          string
}

type FileConfig struct {
	MaxUploadSize int64                  // 10 megabytes = 10 << 20
	CheckTypes    map[string]interface{} // "image/jpeg": nil, "image/png": nil, ...
	Types         []string
}

type Config struct {
	Environment string
	Server      NetServerConfig
	Postgres    database.PostgresConfig
	// DB          db.ConfigDB
	Auth AuthConfig
	File FileConfig
	HTTP HTTPConfig
}
