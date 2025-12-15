package observability

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	promclient "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config groups the knobs used to bootstrap telemetry for the API process.
type Config struct {
	ServiceName  string
	Environment  string
	OTLPEndpoint string
	OTLPHeaders  map[string]string
	MetricsPath  string
}

// Setup configures OpenTelemetry tracing + metrics (Prometheus exporter) and returns
// helpers that can be hooked into routers and shutdown flows.
type Setup struct {
	TracerProvider *sdktrace.TracerProvider
	MeterProvider  *sdkmetric.MeterProvider
	MetricsHandler http.Handler
	GinMiddleware  func() gin.HandlerFunc
}

// New initializes telemetry exporters and providers. It is tolerant of partial
// configuration (e.g. only metrics without OTLP tracing).
func New(ctx context.Context, cfg Config) (*Setup, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource init: %w", err)
	}

	var tp *sdktrace.TracerProvider
	if cfg.OTLPEndpoint != "" {
		clientOpts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.OTLPEndpoint)}
		if strings.HasPrefix(cfg.OTLPEndpoint, "http://") || strings.HasPrefix(cfg.OTLPEndpoint, "localhost:") {
			clientOpts = append(clientOpts, otlptracehttp.WithInsecure())
		}
		if len(cfg.OTLPHeaders) > 0 {
			clientOpts = append(clientOpts, otlptracehttp.WithHeaders(cfg.OTLPHeaders))
		}

		client := otlptracehttp.NewClient(clientOpts...)
		exporter, err := otlptrace.New(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("otlp exporter: %w", err)
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
	}

	registry := promclient.NewRegistry()
	promExporter, err := otelprom.New(otelprom.WithRegisterer(registry))
	if err != nil {
		return nil, fmt.Errorf("prometheus exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	tracerProvider := otel.GetTracerProvider()
	if tp != nil {
		tracerProvider = tp
	}

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	mw := func() gin.HandlerFunc {
		return otelgin.Middleware(cfg.ServiceName,
			otelgin.WithTracerProvider(tracerProvider),
			otelgin.WithMeterProvider(mp),
		)
	}

	return &Setup{
		TracerProvider: tp,
		MeterProvider:  mp,
		MetricsHandler: handler,
		GinMiddleware:  mw,
	}, nil
}

// Shutdown flushes telemetry pipelines best-effort.
func (s *Setup) Shutdown(ctx context.Context) error {
	var errs []error
	if s.TracerProvider != nil {
		errs = append(errs, s.TracerProvider.Shutdown(ctx))
	}
	if s.MeterProvider != nil {
		errs = append(errs, s.MeterProvider.Shutdown(ctx))
	}
	return errors.Join(errs...)
}
