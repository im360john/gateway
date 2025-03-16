package oracle

import (
	"fmt"
	"strings"
)

// Config holds the configuration for Oracle database connection
type Config struct {
	ConnType string   `yaml:"type"`
	Hosts    []string `yaml:"hosts"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
	Schema   string   `yaml:"schema"`
	Port     int      `yaml:"port"`
}

// Type implements the connectors.Config interface
func (c Config) Type() string {
	return "oracle"
}

// Doc returns documentation for the Oracle configuration
func (c Config) Doc() string {
	return "Oracle database connection configuration"
}

// ConnectionString builds Oracle connection string
func (c Config) ConnectionString() string {
	// Oracle connection string format for go-ora:
	// oracle://user:password@host:port/service_name
	host := c.Hosts[0]
	if host == "localhost" {
		host = "127.0.0.1"
	}

	return fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		host,
		c.Port,
		c.Database,
	)
}

// Validate implements the connectors.Config interface
func (c Config) Validate() error {
	var errors []string

	if c.ConnType != "oracle" {
		errors = append(errors, "type must be 'oracle'")
	}
	if len(c.Hosts) == 0 {
		errors = append(errors, "at least one host must be specified")
	}
	if c.User == "" {
		errors = append(errors, "user must be specified")
	}
	if c.Password == "" {
		errors = append(errors, "password must be specified")
	}
	if c.Database == "" {
		errors = append(errors, "database must be specified")
	}
	if c.Schema == "" {
		errors = append(errors, "schema must be specified")
	}
	if c.Port == 0 {
		errors = append(errors, "port must be specified")
	}

	if len(errors) > 0 {
		return fmt.Errorf("invalid oracle configuration: %s", strings.Join(errors, ", "))
	}
	return nil
}

// ExtraPrompt implements the connectors.Config interface
func (c Config) ExtraPrompt() []string {
	return nil
}

// Type implements the connectors.Config interface
func (c Config) GetType() string {
	return c.Type()
}
