package otel

import "time"

// Config represents OpenTelemetry configuration
type Config struct {
	// ExporterType defines the type of exporter to use ("oltp" or "stdout")
	ExporterType string `yaml:"exporter_type"`

	// ServiceName is the name of your service for tracing
	ServiceName string `yaml:"service_name"`

	// ServiceVersion is the version of your service
	ServiceVersion string `yaml:"service_version"`

	// Environment specifies deployment environment (e.g., "prod", "staging")
	Environment string `yaml:"environment"`

	// Endpoint is the OTLP endpoint URL (for oltp exporter)
	Endpoint string `yaml:"endpoint"`

	// Insecure determines if TLS should be disabled
	Insecure bool `yaml:"tls_mode"`

	// SpanMaxQueueSize is the maximum queue size for spans
	SpanMaxQueueSize int `yaml:"span_max_queue_size"`

	// SpanMaxExportBatch is the maximum batch size for span export
	SpanMaxExportBatch int `yaml:"span_max_export_batch"`

	// BatchTimeout is the timeout for batching spans
	BatchTimeout time.Duration `yaml:"batch_timeout"`

	// ResourceAttributes are additional attributes to add to traces
	ResourceAttributes map[string]string `yaml:"resource_attributes"`
}

func (c Config) Tag() string {
	return "otel"
}

func (c Config) WithDefaults() {
	if c.SpanMaxQueueSize == 0 {
		c.SpanMaxQueueSize = 1
	}
	if c.SpanMaxExportBatch == 0 {
		c.SpanMaxExportBatch = 1
	}
	if c.BatchTimeout == 0 {
		c.BatchTimeout = time.Second
	}
}
