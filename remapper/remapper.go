package remapper

import (
	"gopkg.in/yaml.v3"
)

func Remap[TValue any](config any) (TValue, error) {
	var t TValue
	raw, err := yaml.Marshal(config)
	if err != nil {
		return t, err
	}
	if err := yaml.Unmarshal(raw, &t); err != nil {
		return t, err
	}
	return t, nil
}
