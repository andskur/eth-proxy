package middlewares

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"

	"eth-proxy/pkg/logger"
)

// Logger returns a new instance of Logger middleware
// which provides basic request logging
func Logger() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(next, w, r)

			fields := make(logrus.Fields)
			fields["code"] = m.Code
			fields["dt"] = m.Duration
			fields["written"] = m.Written

			var output func(format string, args ...interface{})

			if m.Code == 200 {
				output = logger.Log().WithFields(fields).Infof
			} else {
				output = logger.Log().WithFields(fields).Errorf
			}

			output("%s %s", r.Method, r.URL)
			return
		})
	}
}
