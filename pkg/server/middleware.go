package server

import (
	"dvnetman/pkg/logger"
	"dvnetman/pkg/utils"
	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
	"strings"
)

func isValidRoute(route *mux.Route) bool {
	if route == nil {
		return false
	}
	_, ok := route.GetHandler().(*mux.Router)
	return !ok
}

func logContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			l := logger.Ctx(r.Context()).SubLogger()
			l.Key("remote", r.RemoteAddr).
				Key("method", r.Method).
				Key("url", r.URL.String())
			route := mux.CurrentRoute(r)
			if isValidRoute(route) {
				if route.GetName() != "" {
					l.Key("route", route.GetName())
				}
				if path, err := route.GetPathTemplate(); err == nil && path != "" {
					path = strings.ReplaceAll(
						path, ":"+utils.UUIDRegexString,
						":uuid",
					)
					l.Key("endpoint", path)
				}
			}
			next.ServeHTTP(w, r.WithContext(l.Logger().Context(r.Context())))
		},
	)
}

func accessLogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			m := httpsnoop.CaptureMetrics(handler, w, r)
			logger.Info(r.Context()).
				Key("code", m.Code).
				Key("duration", m.Duration).
				Key("written", m.Written).
				Msg("request")
		},
	)
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(r.Context()).Key("recover", rec).Key("stack", string(debug.Stack())).Msg("Panic recovered")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
