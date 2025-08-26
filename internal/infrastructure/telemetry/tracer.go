package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type Config struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"service_name"`
}

type Tracer struct {
	provider *sdktrace.TracerProvider
}

func NewTracer(ctx context.Context, config Config) (*Tracer, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.URL),
		otlptracehttp.WithInsecure(),
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

	return &Tracer{
		provider: tp,
	}, nil
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}
