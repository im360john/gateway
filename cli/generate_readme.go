package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GenerateReadmeCommand creates a command to generate README.md from CLI commands
func GenerateReadmeCommand() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "generate-docs",
		Short: "Generate CLI documentation",
		Long:  "Generate CLI documentation in Markdown format based on command definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new root command to collect all commands
			rootCmd := &cobra.Command{
				Use:   "gateway",
				Short: "gateway cli",
			}

			// Register all commands to the root command
			RegisterCommand(rootCmd, StartCommand())
			RegisterCommand(rootCmd, Connectors())
			RegisterCommand(rootCmd, Plugins())
			RegisterCommand(rootCmd, Discover())
			RegisterCommand(rootCmd, Connection())
			
			// Add the generate-docs command itself to the documentation
			docCmd := GenerateReadmeCommand()
			RegisterCommand(rootCmd, docCmd)

			// Generate documentation
			doc, err := generateMarkdownDoc(rootCmd)
			if err != nil {
				return err
			}

			// Write to file
			err = os.WriteFile(outputPath, []byte(doc), 0644)
			if err != nil {
				return fmt.Errorf("failed to write documentation to file: %w", err)
			}

			fmt.Printf("Documentation generated successfully at %s\n", outputPath)
			return nil
		},
	}

	// Set default output path to cli/README.md
	defaultPath := filepath.Join("cli", "README.md")
	cmd.Flags().StringVar(&outputPath, "output", defaultPath, "Path to output README.md file")

	return cmd
}

// CommandDocInfo holds information about a command for documentation
type CommandDocInfo struct {
	CommandPath string
	Short       string
	Long        string
	UseLine     string
	Example     string
	HasExample  bool
	Flags       []FlagInfo
	HasFlags    bool
	SubCommands []CommandDocInfo
}

// FlagInfo holds information about a flag for documentation
type FlagInfo struct {
	Name     string
	Usage    string
	DefValue string
}

// generateMarkdownDoc generates markdown documentation for the given command and its subcommands
func generateMarkdownDoc(cmd *cobra.Command) (string, error) {
	tmpl := `---
title: 'Gateway CLI'
---

This document provides information about the available CLI commands and their parameters for the Gateway application.

## Available Commands

{{range .Commands}}### ` + "`{{.CommandPath}}`" + `

{{.Short}}

**Description:**

{{.Long}}

**Usage:**

` + "```" + `
{{.UseLine}}
` + "```" + `

{{if .HasFlags}}**Flags:**

{{range .Flags}}- ` + "`--{{.Name}}`" + ` - {{.Usage}}{{if .DefValue}} (default: "{{.DefValue}}"){{end}}
{{end}}
{{end}}
{{if .HasExample}}
{{.Example}}
{{end}}
{{range .SubCommands}}
### ` + "`{{.CommandPath}}`" + `

{{.Short}}

**Usage:**

` + "```" + `
{{.UseLine}}
` + "```" + `

{{if .HasFlags}}**Flags:**

{{range .Flags}}- ` + "`--{{.Name}}`" + ` - {{.Usage}}{{if .DefValue}} (default: "{{.DefValue}}"){{end}}
{{end}}
{{end}}
{{if .HasExample}}
{{.Example}}
{{end}}
{{end}}
{{end}}

## Configuration File

The gateway.yaml configuration file defines:

- API endpoints
- Database connections
- Security settings
- Plugin configurations

Example configuration:

` + "```yaml" + `
# Example gateway.yaml
api:
  # API configuration
database:
  # Database connection settings
plugins:
  # Plugin configurations
` + "```" + `

For detailed configuration options, please refer to the main documentation.
`

	// Create a template
	t, err := template.New("readme").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data for template
	commands := []CommandDocInfo{}
	for _, command := range cmd.Commands() {
		if command.Hidden {
			continue
		}
		
		commands = append(commands, buildCommandDocInfo(command))
	}

	data := struct {
		Commands []CommandDocInfo
	}{
		Commands: commands,
	}

	// Execute template
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Clean up the output
	output := strings.ReplaceAll(buf.String(), "gateway gateway", "gateway")
	
	return output, nil
}

// buildCommandDocInfo builds documentation info for a command and its subcommands
func buildCommandDocInfo(cmd *cobra.Command) CommandDocInfo {
	info := CommandDocInfo{
		CommandPath: cmd.CommandPath(),
		Short:       cmd.Short,
		Long:        cmd.Long,
		UseLine:     cmd.UseLine(),
		Example:     cmd.Example,
		HasExample:  cmd.Example != "",
		Flags:       []FlagInfo{},
		HasFlags:    false,
		SubCommands: []CommandDocInfo{},
	}

	// Process flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			info.Flags = append(info.Flags, FlagInfo{
				Name:     flag.Name,
				Usage:    flag.Usage,
				DefValue: flag.DefValue,
			})
			info.HasFlags = true
		}
	})

	// Process subcommands
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			info.SubCommands = append(info.SubCommands, buildCommandDocInfo(subCmd))
		}
	}

	return info
} 