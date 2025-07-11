package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type JaegerConfig struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"service_name"`
}

type Jaeger struct {
	provider *sdktrace.TracerProvider
}

func NewJaeger(ctx context.Context, config JaegerConfig) (*Jaeger, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.URL),
		otlptracehttp.WithInsecure(), // Pour le développement, retirez en production si vous utilisez TLS
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)

	return &Jaeger{
		provider: tp,
	}, nil
}

func (t *Jaeger) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
