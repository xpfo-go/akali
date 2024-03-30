package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/url"
	"{{ .ProjectName }}/config"
	"{{ .ProjectName }}/pkg/logs"
	"time"
)

const (
	defaultMaxOpenConn     = 100
	defaultMaxIdleConn     = 25
	defaultConnMaxLifetime = 10 * time.Minute
)

// DBClient MySQL DB Instance
type DBClient struct {
	name            string
	DB              *sqlx.DB
	dataSource      string
	maxOpenConn     int
	maxIdleConn     int
	connMaxLifetime time.Duration
}

// TestConnection ...
func (db *DBClient) TestConnection() (err error) {
	conn, err := sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return
	}

	_ = conn.Close()
	return nil
}

// Connect to db, and update some settings
func (db *DBClient) Connect() error {
	var err error
	db.DB, err = sqlx.Connect("mysql", db.dataSource)
	if err != nil {
		return err
	}

	db.DB.SetMaxOpenConns(db.maxOpenConn)
	db.DB.SetMaxIdleConns(db.maxIdleConn)
	db.DB.SetConnMaxLifetime(db.connMaxLifetime)

	logs.Infof("connect to database: %s[maxOpenConn=%d, maxIdleConn=%d, connMaxLifetime=%s]",
		db.name, db.maxOpenConn, db.maxIdleConn, db.connMaxLifetime)

	return nil
}

// Close db connection
func (db *DBClient) Close() {
	if db.DB != nil {
		_ = db.DB.Close()
	}
}

// NewDBClient :
func NewDBClient(cfg *config.Config) *DBClient {
	dataSource := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&interpolateParams=true&loc=%s&time_zone=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		"utf8",
		"UTC",
		url.QueryEscape("'+00:00'"),
	)

	maxOpenConn := defaultMaxOpenConn
	if cfg.Database.MaxOpenConn > 0 {
		maxOpenConn = cfg.Database.MaxOpenConn
	}

	maxIdleConn := defaultMaxIdleConn
	if cfg.Database.MaxIdleConn > 0 {
		maxIdleConn = cfg.Database.MaxIdleConn
	}

	if maxOpenConn {{ .Lt }} maxIdleConn {
		logs.Errorf("error config for database %s, maxOpenConn should greater or equals to maxIdleConn, will"+
			"use the default [defaultMaxOpenConn=%d, defaultMaxIdleConn=%d]",
			cfg.Database.Name, defaultMaxOpenConn, defaultMaxIdleConn)
		maxOpenConn = defaultMaxOpenConn
		maxIdleConn = defaultMaxIdleConn
	}

	connMaxLifetime := defaultConnMaxLifetime
	if cfg.Database.ConnMaxLifetimeSecond > 0 {
		if cfg.Database.ConnMaxLifetimeSecond >= 60 {
			connMaxLifetime = time.Duration(cfg.Database.ConnMaxLifetimeSecond) * time.Second
		} else {
			logs.Errorf("error config for database %s, connMaxLifetimeSeconds should be greater than 60 seconds"+
				"use the default [defaultConnMaxLifetime=%s]",
				cfg.Database.Name, defaultConnMaxLifetime)
		}
	}

	return &DBClient{
		name:            cfg.Database.Name,
		dataSource:      dataSource,
		maxOpenConn:     maxOpenConn,
		maxIdleConn:     maxIdleConn,
		connMaxLifetime: connMaxLifetime,
	}
}
