package server

import (
	"context"
	"dvnetman/pkg/logger"
	"dvnetman/version"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type otelServer struct {
	otelPropagator    propagation.TextMapPropagator
	otelExporter      sdktrace.SpanExporter
	otelTraceProvider *sdktrace.TracerProvider
	log               logger.Logger
	metricExporter    *prometheus.Exporter
	meterProvider     *metric.MeterProvider
}

func (o *otelServer) setup(ctx context.Context) (err error) {
	o.otelPropagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(o.otelPropagator)
	if o.otelExporter, err = otlptracehttp.New(
		ctx, otlptracehttp.WithEndpoint("localhost:4318"), otlptracehttp.WithInsecure(),
	); err != nil {
		return errors.Wrap(err, "failed to create OTLP sdktrace exporter")
	}
	var res *resource.Resource
	if res, err = resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("dvnetman"),
			semconv.ServiceVersionKey.String(version.Version),
		),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	); err != nil {
		return errors.Wrap(err, "failed to create resource")
	}
	o.otelTraceProvider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(o.otelExporter),
	)
	otel.SetTracerProvider(o.otelTraceProvider)

	if o.metricExporter, err = prometheus.New(); err != nil {
		return errors.Wrap(err, "failed to create Prometheus metric exporter")
	}

	o.meterProvider = metric.NewMeterProvider(
		metric.WithReader(
			o.metricExporter,
		),
	)
	otel.SetMeterProvider(o.meterProvider)

	return nil
}

func (o *otelServer) attach(router *mux.Router) {
	routers := map[*mux.Router]struct{}{}
	_ = router.Walk(
		func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) (_ error) {
			if _, found := routers[router]; !found {
				routers[router] = struct{}{}
				router.Use(o.routeInfo)
			}
			return
		},
	)
	o.log.Debug().Msgf("Attached route info middleware to %v routers", len(routers))
}

func (o *otelServer) routeInfo(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer handler.ServeHTTP(w, r)
			route := mux.CurrentRoute(r)
			if _, ok := route.GetHandler().(*mux.Router); ok {
				return
			}
			span := trace.SpanFromContext(r.Context())
			vars := mux.Vars(r)
			nameSet := false
			if name := route.GetName(); name != "" {
				span.SetName(name)
				span.SetAttributes(semconv.HTTPRouteKey.String(name))
				nameSet = true
			}

			if path, err := route.GetPathTemplate(); err != nil {
				o.log.Error().Msgf("Failed to get path template: %v", err)
			} else {
				span.SetAttributes(semconv.HTTPRouteKey.String(path))
				if !nameSet {
					span.SetName(path)
				}
			}
			for key, value := range vars {
				span.SetAttributes(attribute.String("http.route_param."+key, value))
			}
		},
	)
}

func (o *otelServer) shutdown(ctx context.Context) {
	if err := o.otelTraceProvider.ForceFlush(ctx); err != nil {
		o.log.Error().Msgf("Failed to flush OTLP sdktrace provider: %v", err)
	}
	if err := o.otelExporter.Shutdown(ctx); err != nil {
		o.log.Error().Msgf("Failed to shutdown OTLP sdktrace exporter: %v", err)
	}
	if err := o.otelTraceProvider.Shutdown(ctx); err != nil {
		o.log.Error().Msgf("Failed to shutdown OTLP sdktrace provider: %v", err)
	}
	if err := o.meterProvider.Shutdown(ctx); err != nil {
		o.log.Error().Msgf("Failed to shutdown Prometheus meter provider: %v", err)
	}
}
