package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerProvider struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func NewTracerProvider(ctx context.Context, serviceName, serviceVersion, otlpEndpoint string) (*TracerProvider, error) {
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithInsecure(),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: provider,
		tracer:   provider.Tracer(serviceName),
	}, nil
}

func (tp *TracerProvider) Tracer() trace.Tracer {
	return tp.tracer
}

func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	return tp.provider.Shutdown(ctx)
}

func (tp *TracerProvider) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tp.tracer.Start(ctx, name, opts...)
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, name)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func ContextWithSpan(ctx context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(ctx, span)
}

func AddEvent(ctx context.Context, name string, attrs ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, attrs...)
	}
}

func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err, opts...)
	}
}

func SetSpanStatus(ctx context.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(code, description)
	}
}

func Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

func Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}

func SpanIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().SpanID().String()
	}
	return ""
}
