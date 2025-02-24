package piiremover

import (
	"github.com/centralmind/gateway/plugins"
)

func init() {
	plugins.Register(func(cfg PIIInterceptorConfig) (plugins.Interceptor, error) {
		return New(cfg)
	})
}

type PIIInterceptorConfig struct {
	Columns []string `yaml:"columns" json:"columns"`
}

func (P PIIInterceptorConfig) Tag() string {
	return "pii_remover"
}

type PIIInterceptor struct {
	columns map[string]bool
}

func (p *PIIInterceptor) Doc() string {
	return `
Remove certain column from result
`
}

func (p *PIIInterceptor) Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool) {
	for k := range data {
		if p.columns[k] {
			data[k] = nil
		}
	}
	return data, false
}

func New(config PIIInterceptorConfig) (plugins.Interceptor, error) {
	set := map[string]bool{}
	for _, col := range config.Columns {
		set[col] = true
	}
	return &PIIInterceptor{
		columns: set,
	}, nil
}
