package dbx

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"net/http"
)

const dbExecuterContextKey = "db_executer"

func ToContext(ctx context.Context, stmt Executer) context.Context {
	return context.WithValue(ctx, dbExecuterContextKey, stmt)
}

func ForContext(ctx context.Context) Executer {
	ex, ok := ctx.Value(dbExecuterContextKey).(Executer)
	if !ok || ex == nil {
		log.Error().Msg("Not found Executer in context")
		return nil
	}

	return ex
}

func DbCtxMiddleware(stmt *sqlx.DB) func(http.Handler) http.Handler {
	log.Debug().Msg("Added context db middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ToContext(r.Context(), stmt)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetDbCtx(ctx context.Context) *DbCtx {
	db, isExecuterContext := ForContext(ctx).(ExecuterContext)
	if !isExecuterContext {
		log.Ctx(ctx).Error().Msg("Db executer from context not implemented ExecuterContext interface.")
		return nil
	}
	return &DbCtx{
		db:  db,
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
