package sqlutil

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" //nolint

	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	logger "github.com/meBazil/hr-mvp/pkg/common-logger"
)

const (
	emptyBaseName = ""
)

type ConnectionManager struct {
	cfg             ConnectionManagerConfig
	readConnection  *gorm.DB
	writeConnection *gorm.DB
}

func NewConnectionManager(cfg ConnectionManagerConfig) (*ConnectionManager, error) {
	var err error
	srv := &ConnectionManager{
		cfg: cfg,
	}

	log := logger.GetDefaultLogger()
	gormConfig := &gorm.Config{
		Logger: &Logger{
			Logger:                    &log,
			SlowThreshold:             cfg.SlowThreshold,
			TraceQueries:              cfg.TraceQueries,
			IgnoreRecordNotFoundError: cfg.IgnoreRecordNotFoundError,
		},
	}

	srv.writeConnection, err = gorm.Open(postgres.Open(cfg.WriteDSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't init write repository on %s: %v", cfg.WriteDSN, err)
	}

	writeDB, err := srv.writeConnection.DB()
	if err != nil {
		return nil, fmt.Errorf("couldn't init write connection pool on %s: %v", cfg.WriteDSN, err)
	}

	writeDB.SetMaxOpenConns(cfg.MaxConnections)
	writeDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	if cfg.ReadDSN == "" || cfg.ReadDSN != cfg.WriteDSN {
		return srv, nil
	}

	srv.readConnection, err = gorm.Open(postgres.Open(cfg.ReadDSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't init read repository on %s: %v", cfg.ReadDSN, err)
	}

	readDB, err := srv.readConnection.DB()
	if err != nil {
		return nil, fmt.Errorf("couldn't init read connection pool on %s: %v", cfg.WriteDSN, err)
	}

	readDB.SetMaxOpenConns(cfg.MaxConnections)
	readDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	return srv, nil
}

func (m *ConnectionManager) ApplyMigrations() error {
	master, err := m.GetWriteConnection().DB()
	if err != nil {
		return err
	}

	migrationDriver, err := pgmigrate.WithInstance(master, &pgmigrate.Config{})
	if err != nil {
		return err
	}

	// we use database name from database driver, and setup separate databaseName params as empty
	manager, err := migrate.NewWithDatabaseInstance("file://"+m.cfg.MigrationPath, emptyBaseName, migrationDriver)
	if err != nil {
		return err
	}

	err = manager.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	logger.Info("Migrations were applied")

	return nil
}

func (m *ConnectionManager) GetReadConnection() *gorm.DB {
	if m.readConnection != nil {
		return m.readConnection
	}

	return m.writeConnection
}

func (m *ConnectionManager) GetWriteConnection() *gorm.DB {
	return m.writeConnection
}

func (m *ConnectionManager) GetWriteDB() *sql.DB {
	db, _ := m.GetWriteConnection().DB()

	return db
}

func (m *ConnectionManager) GetReadDB() *sql.DB {
	db, _ := m.GetReadConnection().DB()

	return db
}
