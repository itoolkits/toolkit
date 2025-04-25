// mysql db connection

package mysqldb

import (
	"database/sql"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Option func(db *sql.DB)

// WithMaxOpenConn - config max open conn
func WithMaxOpenConn(maxOpen int) Option {
	return func(db *sql.DB) {
		db.SetMaxOpenConns(maxOpen)
	}
}

// WithMaxIdleConn - config max idle conn
func WithMaxIdleConn(maxIdle int) Option {
	return func(db *sql.DB) {
		db.SetMaxIdleConns(maxIdle)
	}
}

// WithMaxLifeTime - config max lifetime
func WithMaxLifeTime(maxLife time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxLifetime(maxLife)
	}
}

// WithMaxIdleTime - config max idle time
func WithMaxIdleTime(maxIdle time.Duration) Option {
	return func(db *sql.DB) {
		db.SetConnMaxIdleTime(maxIdle)
	}
}

// DB - get mysql db
func DB(dsn string, opts ...Option) (*gorm.DB, error) {
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	myDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(myDB)
	}
	return gdb, nil
}
