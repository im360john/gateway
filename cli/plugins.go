package cli

import (
	"fmt"

	"github.com/centralmind/gateway/plugins"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

// Plugins returns a command that lists all available gateway plugins
// and provides detailed documentation for specific plugins when requested.
func Plugins() *cobra.Command {
	return &cobra.Command{
		Use:   "plugins [plugin-name]",
		Short: "List all available plugins",
		Long: `Display a list of all registered gateway plugins with their configuration documentation.

Plugins extend the functionality of the gateway by adding custom features,
protocols, or integrations. They can be configured in the gateway.yaml file.

When run without arguments, this command lists all available plugins.
When run with a specific plugin name as an argument, it displays detailed
configuration documentation for that plugin.`,
		Example: `  gateway plugins         # List all available plugins
  gateway plugins auth     # Show documentation for the auth plugin
  gateway plugins cache    # Show documentation for the cache plugin`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			knownPlugins := plugins.KnownPlugins()

			// If no specific plugin is requested, list all available plugins
			if len(args) == 0 {
				fmt.Println("Available Gateway Plugins:")

				for _, p := range knownPlugins {
					fmt.Printf("  - %s\n", p.Tag())
				}
				fmt.Println("\nUse 'gateway plugins [plugin-name]' to view detailed documentation for a specific plugin.")
				return nil
			}

			// Display documentation for the requested plugin
			for _, pluginName := range args {
				p, ok := plugins.KnownPlugin(pluginName)
				if !ok {
					fmt.Printf("Error: Plugin '%s' is not available.\n", pluginName)
					fmt.Println("Run 'gateway plugins' to see a list of available plugins.")
					continue
				}
				rawConfig, _ := glamour.RenderWithEnvironmentConfig(p.Doc())
				fmt.Println(rawConfig)
			}
			return nil
		},
	}
}
