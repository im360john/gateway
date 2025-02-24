package otel

import "time"

type Config struct {
	ExporterType       string            `yaml:"exporter_type"`         // Type of exporter, either: oltp or stdout
	Endpoint           string            `yaml:"endpoint"`              // OTLP endpoint (e.g., "localhost:4317")
	ServiceName        string            `yaml:"service_name"`          // Service name
	ServiceVersion     string            `yaml:"service_version"`       // Application version
	Environment        string            `yaml:"environment"`           // Deployment environment (e.g., "dev", "staging", "prod")
	Insecure           bool              `yaml:"insecure"`              // TLS Insecure
	BatchTimeout       time.Duration     `yaml:"batch_timeout"`         // Batch timeout for exporting spans
	SpanMaxQueueSize   int               `yaml:"span_max_queue_size"`   // Max queue size for spans
	SpanMaxExportBatch int               `yaml:"span_max_export_batch"` // Max batch size for exporting spans
	ResourceAttributes map[string]string `yaml:"resource_attributes"`   // Additional resource attributes
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
