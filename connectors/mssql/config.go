package mssql

import (
	_ "embed"
	"fmt"
)

//go:embed readme.md
var docString string

// Config represents the configuration for a Microsoft SQL Server connection
type Config struct {
	Hosts    []string `json:"hosts" yaml:"hosts"`       // List of server addresses
	User     string   `json:"user" yaml:"user"`         // Username for authentication
	Password string   `json:"password" yaml:"password"` // Password for authentication
	Database string   `json:"database" yaml:"database"` // Database name
	Port     int      `json:"port" yaml:"port"`         // Port number (default 1433)
	Schema   string   `json:"schema" yaml:"schema"`     // Schema name (default "dbo")
}

// ConnectionString generates a connection string for Microsoft SQL Server
func (c Config) ConnectionString() string {
	// Use the first host from the list
	server := c.Hosts[0]
	if len(c.Hosts) > 1 {
		// If multiple hosts are provided, we could implement failover here
		// For now, just use the first one
	}

	// Set default schema if not specified
	schema := c.Schema
	if schema == "" {
		schema = "dbo"
	}

	// Build connection string using URL format with port and schema
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&schema=%s",
		c.User,
		c.Password,
		server,
		c.Port,
		c.Database,
		schema,
	)
}

// Validate checks if the configuration is valid
func (c Config) Validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one server address is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.User == "" {
		return fmt.Errorf("username is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	if c.Port == 0 {
		c.Port = 1433 // Set default port if not specified
	}
	return nil
}

// Type returns the type of the connector
func (c Config) Type() string {
	return "mssql"
}

// Doc returns documentation about the configuration
func (c Config) Doc() string {
	return `Microsoft SQL Server connection configuration:
hosts: List of server addresses (e.g., ["localhost", "backup-server"])
user: Username for authentication
password: Password for authentication
database: Database name
port: Port number (default: 1433)
schema: Schema name (default: "dbo")`
}

// ExtraPrompt returns additional prompt information for the configuration
func (c Config) ExtraPrompt() []string {
	return []string{
		"Use symbol ':' instead of '@' for named parameters in sql query",
	}
}
