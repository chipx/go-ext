package dbx

import (
	"context"
	"database/sql"
	"github.com/chipx/go-ext/ctxlog"
	"github.com/jmoiron/sqlx"
)

func GetDbCtx(ctx context.Context) *DbCtx {
	db, isExecuterContext := ForContext(ctx).(ExecuterContext)
	if !isExecuterContext {
		ctxlog.For(ctx).Error("Db executer from context not implemented ExecuterContext interface.")
		return nil
	}
	return &DbCtx{
		db:  db,
		ctx: ctx,
	}
}

func NewDbCtx(ctx context.Context, stmt ExecuterContext) *DbCtx {
	return &DbCtx{
		db:  stmt,
		ctx: ctx,
	}
}

type DbCtx struct {
	db  ExecuterContext
	ctx context.Context
}

func (d *DbCtx) Get(dest interface{}, query string, args ...interface{}) error {
	return sqlx.GetContext(d.ctx, d.db, dest, query, args...)
}

func (d *DbCtx) Select(dest interface{}, query string, args ...interface{}) error {
	return sqlx.SelectContext(d.ctx, d.db, dest, query, args...)
}

func (d *DbCtx) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRowContext(d.ctx, query, args...)
}

func (d *DbCtx) DriverName() string {
	return d.db.DriverName()
}

func (d *DbCtx) Rebind(query string) string {
	return d.db.Rebind(query)
}

func (d *DbCtx) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return d.db.BindNamed(query, arg)
}

func (d *DbCtx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(d.ctx, query, args...)
}

func (d *DbCtx) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return d.db.QueryxContext(d.ctx, query, args...)
}

func (d *DbCtx) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return d.db.QueryRowxContext(d.ctx, query, args...)
}

func (d *DbCtx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(d.ctx, query, args...)
}
