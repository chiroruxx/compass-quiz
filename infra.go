package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

const mysqlNetworkTimeout = 5 * time.Second

var errMissingMySQLConfig = errors.New("MYSQL_ADDR, MYSQL_DATABASE, MYSQL_USER, and MYSQL_PASSWORD are required")

func mysqlConfigFromEnv() *mysql.Config {
	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = os.Getenv("MYSQL_ADDR")
	config.DBName = os.Getenv("MYSQL_DATABASE")
	config.User = os.Getenv("MYSQL_USER")
	config.Passwd = os.Getenv("MYSQL_PASSWORD")
	config.Timeout = mysqlNetworkTimeout
	config.ReadTimeout = mysqlNetworkTimeout
	config.WriteTimeout = mysqlNetworkTimeout

	return config
}

func openMySQL(config *mysql.Config) (*sql.DB, error) {
	if config.Addr == "" || config.DBName == "" || config.User == "" || config.Passwd == "" {
		return nil, errMissingMySQLConfig
	}

	connector, err := mysql.NewConnector(config)
	if err != nil {
		return nil, fmt.Errorf("configure MySQL connection: %w", err)
	}

	return sql.OpenDB(connector), nil
}
