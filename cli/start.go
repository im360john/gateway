package cli

import (
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func StartCommand() *cobra.Command {
	var gatewayParams string
	var addr string
	var servers string
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	}
	cmd.PersistentFlags().StringVar(&gatewayParams, "config", "./gateway.yaml", "path to yaml file with gateway configuration")
	cmd.PersistentFlags().StringVar(&addr, "addr", ":9090", "addr for gateway")
	cmd.PersistentFlags().StringVar(&servers, "servers", "", "comma-separated list of additional server URLs for Swagger UI (e.g., https://dev1.example.com,https://dev2.example.com)")
	RegisterCommand(cmd, REST(&gatewayParams, &addr, &servers))
	RegisterCommand(cmd, MCP(&gatewayParams, &addr))
	RegisterCommand(cmd, MCPStdio(&gatewayParams))
	return cmd
}

// RegisterCommand is like parent.AddCommand(child), but also
// makes chaining of PersistentPreRunE and PersistentPreRun
func RegisterCommand(parent, child *cobra.Command) {
	parentPpre := parent.PersistentPreRunE
	childPpre := child.PersistentPreRunE
	if child.PersistentPreRunE == nil && child.PersistentPreRun != nil {
		childPpre = func(cmd *cobra.Command, args []string) error {
			child.PersistentPreRun(cmd, args)
			return nil
		}
	}
	if childPpre != nil {
		child.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if parentPpre != nil {
				err := parentPpre(cmd, args)
				if err != nil {
					return xerrors.Errorf("cannot process parent PersistentPreRunE: %w", err)
				}
			}
			return childPpre(cmd, args)
		}
	} else if parentPpre != nil {
		child.PersistentPreRunE = parentPpre
	}
	parent.AddCommand(child)
}
