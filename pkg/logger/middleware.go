package logger

import (
	"context"
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

func Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				l := Ctx(r.Context()).SubLogger()
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

func OTelTraceKeyProvider(ctx context.Context) map[string]interface{} {
	s := trace.SpanFromContext(ctx)
	if s.SpanContext().IsValid() {
		return map[string]interface{}{
			"trace": s.SpanContext().TraceID().String(),
			"span":  s.SpanContext().SpanID().String(),
		}
	}
	return nil
}
