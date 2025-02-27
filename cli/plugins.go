package cli

import (
	"fmt"

	"github.com/centralmind/gateway/plugins"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func Plugins() *cobra.Command {
	return &cobra.Command{
		Use:   "plugins",
		Short: "List all available plugins",
		Long:  "Display a list of all registered plugins with their configuration documentation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			knownPlugins := plugins.KnownPlugins()

			if len(args) == 0 {
				fmt.Println("Available Plugins:")

				for _, p := range knownPlugins {
					fmt.Printf("%s\n", p.Tag())
				}
				return nil
			}

			for _, pluginName := range args {
				p, ok := plugins.KnownPlugin(pluginName)
				if !ok {
					fmt.Printf("Plugin: %s is not known\n", pluginName)
					continue
				}
				rawConfig, _ := glamour.RenderWithEnvironmentConfig(p.Doc())
				fmt.Println(rawConfig)
			}
			return nil
		},
	}
}
