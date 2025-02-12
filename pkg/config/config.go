package config

import (
	"encoding/json"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/transformer"
	"gopkg.in/yaml.v3"

	"github.com/doublecloud/transfer/pkg/abstract"
)

type Gateway struct {
	Type         abstract.ProviderType     `yaml:"type"`
	Params       any                       `yaml:"params"`
	Objects      *model.DataObjects        `yaml:"objects"`
	Transformers []transformer.Transformer `yaml:"transformers"`
}

func (g *Gateway) Endpoint() (model.Source, error) {
	return model.NewSource(g.Type, g.ParamRaw())
}

func (g *Gateway) ParamRaw() string {
	switch p := g.Params.(type) {
	case []byte:
		return string(p)
	case string:
		return p
	default:
		data, _ := json.Marshal(p)
		return string(data)
	}
}

func FromYaml(raw []byte) (*Gateway, error) {
	var gw Gateway
	err := yaml.Unmarshal(raw, &gw)
	if err != nil {
		return nil, xerrors.Errorf("unable to parse yaml: %w", err)
	}
	return &gw, nil
}
