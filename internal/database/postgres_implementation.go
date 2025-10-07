package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	db struct {
		connPool *pgx.Conn
	}

	option     func(*string)
	CommandTag = pgconn.CommandTag
	Row        = pgx.Row
	Rows       = pgx.Rows
)

// NewPostgresDB returns a postgres db connection pool
func NewPostgresDB(ctx context.Context, opts ...option) (DB, error) {
	var dsn string
	for _, opt := range opts {
		opt(&dsn)
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return db{connPool: conn}, nil
}

// WithUser adds user
func WithUser(user string) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("user=%s", user)
		} else {
			*s = fmt.Sprintf("%s user=%s", *s, user)
		}
	}
}

// WithPassword adds password
func WithPassword(password string) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("password=%s", password)
		} else {
			*s = fmt.Sprintf("%s password=%s", *s, password)
		}
	}
}

// WithHost adds a host
func WithHost(host string) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("host=%s", host)
		} else {
			*s = fmt.Sprintf("%s host=%s", *s, host)
		}
	}
}

// WithPort adds port
func WithPort(port int) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("port=%d", port)
		} else {
			*s = fmt.Sprintf("%s port=%d", *s, port)
		}
	}
}

// WithDB adds db name
func WithDBName(dbName string) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("dbname=%s", dbName)
		} else {
			*s = fmt.Sprintf("%s dbname=%s", *s, dbName)
		}
	}
}

// WithSslMode sets the ssl mode for db
func WithSslMode(sslMode string) option {
	return func(s *string) {
		if *s == "" {
			*s = fmt.Sprintf("sslmode=%s", sslMode)
		} else {
			*s = fmt.Sprintf("%s sslmode=%s", *s, sslMode)
		}
	}
}

// Ping pings the db server
func (db db) Ping(ctx context.Context) error {
	return db.connPool.Ping(ctx)
}

// Exec executes the query
func (db db) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	return db.connPool.Exec(ctx, query, args...)
}

// Query executes and returns a set or rows
func (db db) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return db.connPool.Query(ctx, query, args...)
}

// QueryRow executes the query and return one row
func (db db) QueryRow(ctx context.Context, query string, args ...any) Row {
	return db.connPool.QueryRow(ctx, query, args...)
}
