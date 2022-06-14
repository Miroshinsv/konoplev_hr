package sqlutil

import (
	"time"
)

type ConnectionManagerConfig struct {
	ReadDSN                   string        `env:"DB_READ_DSN" envDefault:""`
	WriteDSN                  string        `env:"DB_WRITE_DSN" envDefault:""`
	MigrationPath             string        `env:"DB_MIGRATION_PATH" envDefault:"/opt/migrations"`
	SlowThreshold             time.Duration `env:"DB_SLOW_THRESHOLD" envDefault:"200ms"`
	TraceQueries              bool          `env:"DB_TRACE_QUERIES" envDefault:"false"`
	IgnoreRecordNotFoundError bool          `env:"DB_IGNORE_RECORD_NOT_FOUND_ERROR" envDefault:"false"`
	MaxConnections            int           `env:"DB_MAX_CONNECTIONS" envDefault:"30"`
	MaxIdleConnections        int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	StatsTickerSeconds        int           `env:"DB_STATS_TICKER_SECONDS" envDefault:"15"`
}
