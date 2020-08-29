package dbx

import (
	"context"
	log "github.com/sirupsen/logrus"
)

const dbExecuterContextKey = "db_executer"

func ToContext(ctx context.Context, stmt Executer) context.Context {
	return context.WithValue(ctx, dbExecuterContextKey, stmt)
}

func ForContext(ctx context.Context) Executer {
	ex, ok := ctx.Value(dbExecuterContextKey).(Executer)
	if !ok || ex == nil {
		log.Error("Not found Executer in context")
		return nil
	}

	return ex
}
