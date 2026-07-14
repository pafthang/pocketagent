package telemetry

import (
	"context"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	once     sync.Once
	enabled  bool
	shutdown func(context.Context) error
)

// Init configures OTLP tracing when OTEL_EXPORTER_OTLP_ENDPOINT is set.
func Init(serviceName string) {
	once.Do(func() {
		if os.Getenv("OTEL_SDK_DISABLED") == "true" || os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
			return
		}

		ctx := context.Background()
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			return
		}

		svcName := serviceName
		if override := os.Getenv("OTEL_SERVICE_NAME"); override != "" {
			svcName = override
		}

		res, err := resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(svcName),
			),
		)
		if err != nil {
			_ = exporter.Shutdown(ctx)
			return
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		enabled = true
		shutdown = tp.Shutdown
	})
}

// Enabled reports whether OTLP tracing is configured via env.
func Enabled() bool {
	return os.Getenv("OTEL_SDK_DISABLED") != "true" && os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != ""
}

// EchoMiddleware traces HTTP requests when OTLP is configured.
func EchoMiddleware(serviceName string) echo.MiddlewareFunc {
	Init(serviceName)
	if !enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}
	return otelecho.Middleware(serviceName)
}

// Shutdown flushes pending spans.
func Shutdown(ctx context.Context) error {
	if shutdown == nil {
		return nil
	}
	return shutdown(ctx)
}