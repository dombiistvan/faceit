package db

import (
	"faceit/db/models"
	"fmt"
	"strings"
	"time"
)

var (
	dbManager *Manager
)

// GetManager returns database manager
func GetManager() *Manager {
	return dbManager
}

// SetupDatabase sets up database connection based on config
func SetupDatabase(dbConf Config) error {
	if dbManager != nil {
		return fmt.Errorf("connection has already been estbilished")
	}

	err := dbConf.Validate()
	if err != nil {
		return err
	}

	dbManager, err = NewManager(&dbConf)
	if err != nil {
		return err
	}

	return dbManager.Conn(nil).Migrator().AutoMigrate(
		&models.User{},
	)
}

// NewManager initialize manager
func NewManager(config *Config) (*Manager, error) {
	var err error

	if config == nil {
		return nil, fmt.Errorf("database config is nil")
	}

	if err = config.Validate(); err != nil {
		return nil, fmt.Errorf("database config is invalid: %w", err)
	}

	// try to connect to DATABASE
	gormDBConn, err := getGormConnection(*config, true)

	if err != nil {
		if !strings.Contains(err.Error(), "Unknown database") {
			return nil, fmt.Errorf("could not estabilish database connection: %w", err)
		}

		// DATABASE HAS NOT FOUND, connect to server without DB
		gormDBConn, err = getGormConnection(*config, false)
		if err != nil {
			return nil, fmt.Errorf("could not estabilish server connection: %w", err)
		}

		// tries to connect to DB
		gormDBConn, err = getGormConnection(*config, true)
		if err != nil {
			// try to create database
			if err = gormDBConn.Exec(fmt.Sprintf("CREATE DATABASE %s collate utf8mb4_general_ci;", config.DBName)).Error; err != nil {
				return nil, fmt.Errorf("could not create database: %w", err)
			}

			gormDBConn, err = getGormConnection(*config, true)
			if err != nil {
				return nil, fmt.Errorf("could not estabilish database connection after create: %w", err)
			}
		}
	}

	gormDBConn.Set("gorm:table_options", "ENGINE=InnoDB")

	sqlDB, err := gormDBConn.DB()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve sql conn for gorm db: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Manager{
		gorm: gormDBConn,
	}, nil
}
