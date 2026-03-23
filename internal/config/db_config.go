package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	dbUserEnvName = "DB_USER"
	dbPassEnvName = "DB_PASSWORD"
	dbNameEnvName = "DATABASE_NAME"
	dbHostEnvName = "DB_HOST"
	dbPortEnvName = "DB_PORT"
)

type DBConfig interface {
	GetDSN() string
}

type dbConfig struct {
	user     string
	password string
	dbName   string
	host     string
	port     string
}

func NewDBConfig() (DBConfig, error) {
	user := os.Getenv(dbUserEnvName)
	if len(user) == 0 {
		return nil, errors.New(dbUserEnvName + " is not set")
	}
	password := os.Getenv(dbPassEnvName)
	if len(password) == 0 {
		return nil, errors.New(dbPassEnvName + " is not set")
	}
	dbName := os.Getenv(dbNameEnvName)
	if len(dbName) == 0 {
		return nil, errors.New(dbNameEnvName + " is not set")
	}
	host := os.Getenv(dbHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(dbHostEnvName + " is not set")
	}
	port := os.Getenv(dbPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(dbPortEnvName + " is not set")
	}
	return &dbConfig{
		host:     host,
		port:     port,
		dbName:   dbName,
		user:     user,
		password: password,
	}, nil
}
func (dbc *dbConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbc.host, dbc.port, dbc.dbName, dbc.user, dbc.password)
}
