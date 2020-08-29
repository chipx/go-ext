package dbx

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Executer interface {
	sqlx.Ext
	//
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	QueryRow(query string, args ...interface{}) *sql.Row
}

type ExecuterContext interface {
	sqlx.ExtContext

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
