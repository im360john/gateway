package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/centralmind/gateway/cli"
	_ "github.com/centralmind/gateway/connectors/postgres"
	_ "github.com/centralmind/gateway/plugins/lua_rls"
	_ "github.com/centralmind/gateway/plugins/pii_remover"
)

func main() {
	rootCommand := &cobra.Command{
		Use:          "gateway",
		Short:        "gateway cli",
		Example:      "./gateway help",
		SilenceUsage: true,
	}
	cli.RegisterCommand(rootCommand, cli.StartCommand())
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
