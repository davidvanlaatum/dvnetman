package logger

import (
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func isValidRoute(route *mux.Route) bool {
	if route == nil {
		return false
	}
	_, ok := route.GetHandler().(*mux.Router)
	return !ok
}

func Middleware(log Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				t := trace.SpanFromContext(r.Context())
				l := log.SubLogger()
				if t.SpanContext().IsValid() {
					l.Key("trace", t.SpanContext().TraceID().String()).
						Key("span", t.SpanContext().SpanID().String())
				}
				l.Key("remote", r.RemoteAddr).
					Key("method", r.Method).
					Key("url", r.URL.String())
				route := mux.CurrentRoute(r)
				if isValidRoute(route) {
					if route.GetName() != "" {
						l.Key("route", route.GetName())
					}
					if path, err := route.GetPathTemplate(); err == nil && path != "" {
						l.Key("endpoint", path)
					}
				}
				r = r.WithContext(l.Logger().Context(r.Context()))
				next.ServeHTTP(w, r)
			},
		)
	}
}
