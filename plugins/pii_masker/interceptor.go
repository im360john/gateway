package piimasker

import "github.com/centralmind/gateway/plugins"

type PIIInterceptorConfig struct {
	Columns []string `yaml:"columns" json:"columns"`
}

type PIIInterceptor struct {
	columns map[string]bool
}

func (p *PIIInterceptor) Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool) {
	for k := range data {
		if p.columns[k] {
			data[k] = nil
		}
	}
	return procesed, false
}

func init() {
	plugins.RegisterInterceptor("pii_remover", func(cfg any) (plugins.Interceptor, error) {
		return New(cfg.(PIIInterceptorConfig))
	})
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
