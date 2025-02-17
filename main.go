package main

import (
	"github.com/centralmind/gateway/pkg/cli"
	_ "github.com/doublecloud/transfer/pkg/transformer/registry"
	"os"

	"github.com/doublecloud/transfer/pkg/cobraaux"
	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{
		Use:          "gateway",
		Short:        "gateway cli",
		Example:      "./gateway help",
		SilenceUsage: true,
	}
	cobraaux.RegisterCommand(rootCommand, cli.StartCommand())
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}
