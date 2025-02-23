package plugins

import (
	"github.com/centralmind/gateway/remapper"
)

type Config interface {
	Tag() string
}

type Interceptor interface {
	Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool)
}

var interceptors = map[string]func(any) (Interceptor, error){}

func RegisterInterceptor[TConfig Config](f func(cfg TConfig) (Interceptor, error)) {
	var t TConfig
	interceptors[t.Tag()] = func(a any) (Interceptor, error) {
		cfg, err := remapper.Remap[TConfig](a)
		if err != nil {
			return nil, xerrors.Errorf("unable to rempa: %w", err)
		}
		return f(cfg)
	}
}

func New(tag string, config any) (Interceptor, error) {
	f, ok := interceptors[tag]
	if !ok {
		return nil, xerrors.Errorf("plugin: %s not found", tag)
	}
	return f(config)
}
