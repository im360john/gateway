package piiremover

import (
	"github.com/centralmind/gateway/plugins"
)

func init() {
	plugins.Register(func(cfg Config) (plugins.Interceptor, error) {
		return New(cfg)
	})
}

type Plugin struct {
	columns map[string]bool
	cfg     Config
}

func (p *Plugin) Doc() string {
	return `
Remove certain column from result

# Example YAML configuration:

pii_remover:
  fields:
    - "*.email"
    - "users.phone"
    - "*.credit_card"
  replacement: "*[REDACTED]*"
  detection_rules:
    credit_card: "\\d{4}-\\d{4}-\\d{4}-\\d{4}"
    phone: "\\+?\\d{10,12}"
`
}

func (p *Plugin) Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool) {
	for k := range data {
		if p.columns[k] {
			data[k] = p.cfg.Replacement
		}
	}
	return data, false
}

func New(config Config) (plugins.Interceptor, error) {
	set := map[string]bool{}
	for _, col := range config.Fields {
		set[col] = true
	}
	return &Plugin{
		columns: set,
		cfg:     config,
	}, nil
}
