package cli

import (
	"fmt"

	"github.com/centralmind/gateway/connectors"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

// Connectors returns a command that lists all available database connectors
// and provides detailed documentation for specific connectors when requested.
func Connectors() *cobra.Command {
	return &cobra.Command{
		Use:   "connectors [connector-name]",
		Short: "List all available database connectors",
		Long: `Display a list of all registered database connectors with their configuration documentation.

When run without arguments, this command lists all available database connectors.
When run with a specific connector name as an argument, it displays detailed
configuration documentation for that connector.

Examples:
  gateway connectors         # List all available connectors
  gateway connectors postgres # Show documentation for PostgreSQL connector
  gateway connectors mysql    # Show documentation for MySQL connector`,
		Example: `  gateway connectors
  gateway connectors postgres
  gateway connectors mysql`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			knownConnectors := connectors.KnownConnectors()

			// If no specific connector is requested, list all available connectors
			if len(args) == 0 {
				fmt.Println("Available Database Connectors:")

				for _, c := range knownConnectors {
					fmt.Printf("  - %s\n", c.Type())
				}
				fmt.Println("\nUse 'gateway connectors [connector-name]' to view detailed documentation for a specific connector.")
				return nil
			}
			
			// Display documentation for the requested connector
			for _, connector := range args {
				c, ok := connectors.KnownConnector(connector)
				if !ok {
					fmt.Printf("Error: Connector '%s' is not available.\n", connector)
					fmt.Println("Run 'gateway connectors' to see a list of available connectors.")
					return nil
				}
				rawConfig, _ := glamour.RenderWithEnvironmentConfig(c.Doc())
				fmt.Println(rawConfig)
			}
			return nil
		},
	}
}
