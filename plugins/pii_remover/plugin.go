package piiremover

import (
	"fmt"
	"path"
	"regexp"

	"github.com/centralmind/gateway/plugins"
)

func init() {
	plugins.Register(func(cfg Config) (plugins.Interceptor, error) {
		return New(cfg)
	})
}

type Plugin struct {
	patterns map[string]*regexp.Regexp
	columns  map[string]bool
	cfg      Config
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

func (p *Plugin) Process(data map[string]any, context map[string][]string) (processed map[string]any, skipped bool) {
	for k, v := range data {
		if p.columns[k] {
			data[k] = p.cfg.Replacement
			continue
		}

		for pattern := range p.columns {
			if matched, _ := path.Match(pattern, k); matched {
				data[k] = p.cfg.Replacement
				break
			}
		}

		if strVal, ok := v.(string); ok {
			for name, regex := range p.patterns {
				if k != name {
					continue
				}
				if regex.MatchString(strVal) {
					data[k] = p.cfg.Replacement
					break
				}
			}
		}
	}
	return data, false
}

func New(config Config) (plugins.Interceptor, error) {
	p := &Plugin{
		patterns: make(map[string]*regexp.Regexp),
		columns:  make(map[string]bool),
		cfg:      config,
	}

	for name, pattern := range config.DetectionRules {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid detection rule pattern for %s: %v", name, err)
		}
		p.patterns[name] = regex
	}

	if p.cfg.Replacement == "" {
		p.cfg.Replacement = "[REDACTED]"
	}

	for _, field := range config.Fields {
		p.columns[field] = true
	}

	return p, nil
}
