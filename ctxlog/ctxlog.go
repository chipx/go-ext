package ctxlog

import "context"
import log "github.com/sirupsen/logrus"

const ctxLogFieldsKey = "ctx_log_fields"

type ctxLogFields struct {
	log.Fields
}

func For(ctx context.Context) log.FieldLogger {
	fields, ok := ctx.Value(ctxLogFieldsKey).(*ctxLogFields)
	if !ok || fields == nil {
		return log.StandardLogger()
	}

	return log.WithFields(fields.Fields)
}

func To(ctx context.Context, fields log.Fields) context.Context {
	l, ok := ctx.Value(ctxLogFieldsKey).(*ctxLogFields)
	if !ok || l == nil {
		l = &ctxLogFields{fields}
		ctx = context.WithValue(ctx, ctxLogFieldsKey, l)
	}

	for k, v := range fields {
		l.Fields[k] = v
	}

	return ctx
}

func AddFields(ctx context.Context, fields log.Fields) {
	l, ok := ctx.Value(ctxLogFieldsKey).(*ctxLogFields)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.Fields[k] = v
	}
}
