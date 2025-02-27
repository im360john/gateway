package model

import (
	"encoding/json"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	API      APIParams      `yaml:"api" json:"api"`
	Database Database       `yaml:"database" json:"database"`
	Plugins  map[string]any `yaml:"plugins" json:"plugins"`
}

func FromYaml(raw []byte) (*Config, error) {
	var gw Config
	err := yaml.Unmarshal(raw, &gw)
	if err != nil {
		return nil, xerrors.Errorf("unable to parse yaml: %w", err)
	}
	return &gw, nil
}

func (g *Config) ParamRaw() string {
	switch p := g.Database.Connection.(type) {
	case []byte:
		return string(p)
	case string:
		return p
	default:
		data, _ := json.Marshal(p)
		return string(data)
	}
}

type APIParams struct {
	Name        string `yaml:"name" json:"name,omitempty"`
	Description string `yaml:"description" json:"description,omitempty"`
	Version     string `yaml:"version" json:"version,omitempty"`
}

type Database struct {
	Type       string  `yaml:"type" json:"type,omitempty"`
	Connection any     `yaml:"connection" json:"connection,omitempty"`
	Tables     []Table `yaml:"tables" json:"tables,omitempty"`
}

type Table struct {
	Name      string         `yaml:"name" json:"name,omitempty"`
	Columns   []ColumnSchema `yaml:"columns" json:"columns,omitempty"`
	Endpoints []Endpoint     `yaml:"endpoints" json:"endpoints,omitempty"`
}

type ColumnSchema struct {
	Name       string `yaml:"name" json:"name,omitempty"`
	Type       string `yaml:"type" json:"type,omitempty"`
	PrimaryKey bool   `yaml:"primary_key" json:"primary_key,omitempty"`
	PII        bool   `yaml:"pii" json:"pii,omitempty"`
}

type Endpoint struct {
	HTTPMethod  string           `yaml:"http_method" json:"http_method,omitempty"`
	HTTPPath    string           `yaml:"http_path" json:"path,omitempty"`
	MCPMethod   string           `yaml:"mcp_method" json:"mcp_method,omitempty"`
	Summary     string           `yaml:"summary" json:"summary,omitempty"`
	Description string           `yaml:"description" json:"description,omitempty"`
	Query       string           `yaml:"query" json:"query,omitempty"`
	Params      []EndpointParams `yaml:"params" json:"params,omitempty"`
}

type EndpointParams struct {
	Name     string      `yaml:"name" json:"name,omitempty"`
	Type     string      `yaml:"type" json:"type,omitempty"`
	Location string      `yaml:"location" json:"location,omitempty"`
	Required bool        `yaml:"required" json:"required,omitempty"`
	Format   string      `yaml:"format,omitempty" json:"format,omitempty"`
	Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`
}
