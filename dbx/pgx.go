package dbx

import (
	"context"
	"github.com/chipx/go-ext/ctxlog"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func NewPgxLogger(fields log.Fields) pgx.Logger {
	return &PgxLogger{fields: fields}
}

type PgxLogger struct {
	fields log.Fields
}

func (l PgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	logger := ctxlog.For(ctx).WithFields(l.fields)
	if data != nil {
		logger = logger.WithFields(data)
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.WithField("PGX_LOG_LEVEL", level).Debug(msg)
	case pgx.LogLevelDebug:
		logger.Debug(msg)
	case pgx.LogLevelInfo:
		logger.Info(msg)
	case pgx.LogLevelWarn:
		logger.Warn(msg)
	case pgx.LogLevelError:
		logger.Error(msg)
	default:
		logger.WithField("INVALID_PGX_LOG_LEVEL", level).Error(msg)
	}
}
