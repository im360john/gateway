package cli

import (
	"fmt"

	"github.com/centralmind/gateway/connectors"
	"github.com/spf13/cobra"
)

func Connectors() *cobra.Command {
	return &cobra.Command{
		Use:   "connectors",
		Short: "List all available database connectors",
		Long:  "Display a list of all registered database connectors with their configuration documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			knownConnectors := connectors.KnownConnectors()
			
			fmt.Println("Available Database Connectors:")
			fmt.Println("=============================")
			
			for _, c := range knownConnectors {
				fmt.Printf("\nConnector: %s\n", c.Type())
				fmt.Printf("--------------%s\n", string(make([]byte, len(c.Type()))))
				fmt.Println(c.Doc())
			}
			
			return nil
		},
	}
} 