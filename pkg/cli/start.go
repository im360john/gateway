package cli

import (
	"github.com/doublecloud/transfer/pkg/cobraaux"
	"github.com/spf13/cobra"
)

func StartCommand() *cobra.Command {
	var gatewayParams string
	var addr string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	}
	cmd.PersistentFlags().StringVar(&gatewayParams, "config", "./gateway.yaml", "path to yaml file with gateway configuration")
	cmd.PersistentFlags().StringVar(&addr, "addr", ":9090", "addr for gateway")
	cobraaux.RegisterCommand(cmd, REST(&gatewayParams, &addr))
	cobraaux.RegisterCommand(cmd, MCP(&gatewayParams, &addr))
	cobraaux.RegisterCommand(cmd, MCPStdio(&gatewayParams))
	return cmd
}
