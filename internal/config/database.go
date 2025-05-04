package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type MysqlDatabaseConfig struct {
	Hostname              string
	Port                  int
	Username              string
	Password              string
	Database              string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifeTime time.Duration
	ConnectionMaxIdleTime time.Duration
	ParseTime             bool
}

func (m *MysqlDatabaseConfig) GetDsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=%t",
		m.Username,
		m.Password,
		m.Hostname,
		m.Port,
		m.Database,
		m.ParseTime,
	)
}

func NewMysqlDatabaseConfig() *MysqlDatabaseConfig {
	config := MysqlDatabaseConfig{
		Hostname:              os.Getenv("MYSQL_HOST"),
		Port:                  3306,
		Username:              os.Getenv("MYSQL_USER"),
		Password:              os.Getenv("MYSQL_PASSWORD"),
		Database:              os.Getenv("MYSQL_DATABASE"),
		MaxOpenConnections:    10,
		MaxIdleConnections:    10,
		ConnectionMaxLifeTime: 6 * time.Minute,
		ConnectionMaxIdleTime: 6 * time.Minute,
		ParseTime:             true,
	}

	if port, err := strconv.Atoi(os.Getenv("MYSQL_PORT")); err == nil && port > 0 {
		config.Port = port
	}

	if MaxOpenConn, err := strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN_CONNECTIONS")); err == nil && MaxOpenConn > 0 {
		config.MaxOpenConnections = MaxOpenConn
	}

	if MaxIdleConn, err := strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE_CONNECTIONS")); err == nil && MaxIdleConn > 0 {
		config.MaxOpenConnections = MaxIdleConn
	}

	if MaxLifeTime, err := strconv.Atoi(os.Getenv("MYSQL_CONNECTION_MAX_LIFE_TIME")); err == nil && MaxLifeTime > 0 {
		config.ConnectionMaxLifeTime = time.Duration(MaxLifeTime) * time.Second
	}

	if MaxIdleTime, err := strconv.Atoi(os.Getenv("MYSQL_CONNECTION_MAX_IDLE_TIME")); err == nil && MaxIdleTime > 0 {
		config.ConnectionMaxIdleTime = time.Duration(MaxIdleTime) * time.Second
	}

	if "false" == os.Getenv("MYSQL_PARSE_TIME") {
		config.ParseTime = false
	}

	return &config
}
