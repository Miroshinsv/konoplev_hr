package config

import (
	"github.com/meBazil/hr-mvp/pkg/sqlutil"
)

type App struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	Health   Health

	REST REST
	DB   sqlutil.ConnectionManagerConfig
	JWT  JWT
	ACL  ACL
}
