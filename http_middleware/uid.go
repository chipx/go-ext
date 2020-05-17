package http_middleware

import (
	"github.com/chipx/go-ext/ctxlog"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func UidMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startRequest := time.Now()

			var uuidStr string
			if u, err := uuid.NewV4(); err == nil {
				uuidStr = u.String()
			} else {
				log.WithError(err).Error("Fail generate uuid")
				uuidStr = strconv.Itoa(int(time.Now().UnixNano()))
			}

			ctx := ctxlog.To(r.Context(), log.Fields{"uid": uuidStr})

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			ctxlog.For(ctx).Debugf("Request executed at %s", time.Since(startRequest))
		})
	}
}
