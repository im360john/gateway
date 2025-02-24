package otel

import (
	"context"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	trace_provider "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func init() {
	plugins.Register[Config](New)
}

type Plugin struct {
	config Config
	tp     *trace_provider.TracerProvider
}

func (p Plugin) Wrap(connector connectors.Connector) (connectors.Connector, error) {
	return Connector{
		inner:  connector,
		config: p.config,
		tp:     p.tp,
	}, nil
}

func (p Plugin) Doc() string {
	return `
Allow to configure otel exporter, example config:

    exporter_type: oltp
    service_name: gachi_bass
    endpoint: localhost:4317
    tls_mode: insecure
    span_max_queue_size: 5
    span_max_export_batch: 10
    batch_timeout: 1s
`
}

func New(cfg Config) (plugins.Wrapper, error) {
	ctx := context.Background()
	var err error
	var exporter trace_provider.SpanExporter
	cfg.WithDefaults()
	exporter, err = stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}
	switch cfg.ExporterType {
	case "oltp":
		var opts []otlptracegrpc.Option
		opts = append(
			opts,
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
			otlptracegrpc.WithDialOption(grpc.WithBlock()),
		)
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
		if err != nil {
			return nil, err
		}
	default:
	}

	// Define resource attributes (service, version, env)
	resAttrs := []attribute.KeyValue{
		attribute.String("service.name", cfg.ServiceName),
		attribute.String("service.version", cfg.ServiceVersion),
		attribute.String("deployment.environment", cfg.Environment),
	}

	// Add custom resource attributes if provided
	for key, value := range cfg.ResourceAttributes {
		resAttrs = append(resAttrs, attribute.String(key, value))
	}

	// Create a new resource with attributes
	res, err := resource.New(ctx, resource.WithAttributes(resAttrs...))
	if err != nil {
		return nil, err
	}

	// Create a new tracer provider with batch settings
	tp := trace_provider.NewTracerProvider(
		trace_provider.WithBatcher(exporter,
			trace_provider.WithBatchTimeout(cfg.BatchTimeout),
			trace_provider.WithMaxQueueSize(cfg.SpanMaxQueueSize),
			trace_provider.WithMaxExportBatchSize(cfg.SpanMaxExportBatch),
		),
		trace_provider.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)
	return &Plugin{
		config: cfg,
		tp:     tp,
	}, nil
}
