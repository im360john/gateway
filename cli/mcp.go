package cli

import (
	"context"
	"os"

	"github.com/centralmind/gateway/mcpgenerator"
	gw_model "github.com/centralmind/gateway/model"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func MCP(configPath *string, addr *string) *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "MCP gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) error {
			gwRaw, err := os.ReadFile(*configPath)
			if err != nil {
				return xerrors.Errorf("unable to read yaml config file: %w", err)
			}
			gw, err := gw_model.FromYaml(gwRaw)
			if err != nil {
				return xerrors.Errorf("unable to parse config file: %w", err)
			}
			srv, err := mcpgenerator.New(*gw)
			if err != nil {
				return xerrors.Errorf("unable to init mcp generator: %w", err)
			}
			return srv.ServeSSE(*addr).Start(*addr)
		},
	}
}

func MCPStdio(configPath *string) *cobra.Command {
	var logFile string
	res := &cobra.Command{
		Use:   "mcp-stdio",
		Short: "MCP gateway via std-io",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) error {
			gwRaw, err := os.ReadFile(*configPath)
			if err != nil {
				return xerrors.Errorf("unable to read yaml config file: %w", err)
			}
			gw, err := gw_model.FromYaml(gwRaw)
			if err != nil {
				return xerrors.Errorf("unable to parse config file: %w", err)
			}
			srv, err := mcpgenerator.New(*gw)
			if err != nil {
				return xerrors.Errorf("unable to init mcp generator: %w", err)
			}
			return srv.ServeStdio().Listen(context.Background(), os.Stdin, os.Stdout)
		},
	}
	res.Flags().StringVar(&logFile, "log-file", "/var/log/gateway/mcp.log", "path to log file")
	return res
}
