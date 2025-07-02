package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
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

func NewJaeger(config JaegerConfig) (*Jaeger, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.URL)))
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
