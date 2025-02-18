package plugins

import "github.com/pkg/errors"

type Interceptor interface {
	Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool)
}

var interceptors = map[string]func(any) (Interceptor, error){}

func RegisterInterceptor(tag string, f func(cfg any) (Interceptor, error)) {
	interceptors[tag] = f
}

func New(tag string, config any) (Interceptor, error) {
	f, ok := interceptors[tag]
	if !ok {
		return nil, errors.Errorf("plugin: %s not found", tag)
	}
	return f(config)
}
