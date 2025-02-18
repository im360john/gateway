package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/centralmind/gateway/pkg/cli"
	_ "github.com/centralmind/gateway/pkg/plugins/lua_rls"
	_ "github.com/centralmind/gateway/pkg/plugins/pii_masker"
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
