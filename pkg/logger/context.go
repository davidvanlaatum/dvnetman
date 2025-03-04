package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

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
