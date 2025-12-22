package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	URL         string `mapstructure:"url"`
	ServiceName string `mapstructure:"service_name"`
}

type Telemetry struct {
	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider
}

func NewTelemetry(ctx context.Context, config Config) (*Telemetry, error) {
	var err error
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(config.ServiceName),
	)

	tp, err := NewTracer(ctx, config, res)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider: %w", err)
	}

	mp, err := NewMeter(ctx, config, res)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider: %w", err)
	}

	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(mp)

	return &Telemetry{
		tracerProvider: tp,
		meterProvider:  mp,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	errMeter := t.meterProvider.Shutdown(ctx)
	errTracer := t.tracerProvider.Shutdown(ctx)

	if errMeter != nil && errTracer != nil {
		return fmt.Errorf("failed to shutdown tracer and meter: %w, %w", errMeter, errTracer)
	}

	if errMeter != nil {
		return fmt.Errorf("failed to shutdown meter: %w", errMeter)
	}

	if errTracer != nil {
		return fmt.Errorf("failed to shutdown tracer: %w", errTracer)
	}

	return nil
}
