package ctxlog

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HttpRequestLogMiddleware(requestIdToResponse bool) func(next http.Handler) http.Handler {
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

			ctx := To(r.Context(), log.Fields{
				"request-id": uuidStr,
				"remote-ip":  strings.Split(r.RemoteAddr, ":")[0],
				"path":       r.URL.Path,
			})

			if requestIdToResponse {
				w.Header().Add("Request-Id", uuidStr)
			}

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			For(ctx).Debugf("Request executed at %s", time.Since(startRequest))
		})
	}
}
