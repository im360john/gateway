package cli

import (
	"fmt"

	"github.com/centralmind/gateway/connectors"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func Connectors() *cobra.Command {
	return &cobra.Command{
		Use:   "connectors",
		Short: "List all available database connectors",
		Long:  "Display a list of all registered database connectors with their configuration documentation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			knownConnectors := connectors.KnownConnectors()

			if len(args) == 0 {
				fmt.Println("Available Database Connectors:")

				for _, c := range knownConnectors {
					fmt.Printf("%s\n", c.Type())
				}
			}
			for _, connector := range args {
				c, ok := connectors.KnownConnector(connector)
				if !ok {
					fmt.Printf("Connector: %s is not known\n", connector)
					return nil
				}
				fmt.Printf("Connector: %s\n", c.Type())
				fmt.Printf("--------------\n")
				rawConfig, _ := glamour.RenderWithEnvironmentConfig(c.Doc())
				fmt.Println(rawConfig)
			}
			return nil
		},
	}
}
