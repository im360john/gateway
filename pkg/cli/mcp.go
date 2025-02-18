package cli

import (
	"context"
	"github.com/centralmind/gateway/pkg/logger"
	"github.com/centralmind/gateway/pkg/mcpgenerator"
	"github.com/doublecloud/transfer/library/go/core/metrics/solomon"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract/coordinator"
	"github.com/doublecloud/transfer/pkg/providers"
	"github.com/doublecloud/transfer/pkg/runtime/local"
	"github.com/spf13/cobra"
	"os"
)

func MCP(configPath *string, addr *string) *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "MCP gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) error {
			transfer, _, res, err := PrepareConfig(configPath)
			if err != nil {
				return xerrors.Errorf("unable to prepare config: %w", err)
			}
			sF, ok := providers.Source[providers.Snapshot](
				logger.NewConsoleLogger(),
				solomon.NewRegistry(solomon.NewRegistryOpts()),
				coordinator.NewFakeClient(),
				transfer,
			)
			if !ok {
				return xerrors.Errorf("no snapshot provider: %T", transfer.Src)
			}
			srv, err := mcpgenerator.New(res.Schema, res.Preview, sF, transfer)
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
			local.WithLogger(logger.NewFileLog(logFile))
			transfer, _, res, err := PrepareConfig(configPath)
			if err != nil {
				return xerrors.Errorf("unable to prepare config: %w", err)
			}
			sF, ok := providers.Source[providers.Snapshot](
				logger.NewFileLog(logFile),
				solomon.NewRegistry(solomon.NewRegistryOpts()),
				coordinator.NewFakeClient(),
				transfer,
			)
			if !ok {
				return xerrors.Errorf("no snapshot provider: %T", transfer.Src)
			}
			srv, err := mcpgenerator.New(res.Schema, res.Preview, sF, transfer)
			if err != nil {
				return xerrors.Errorf("unable to init mcp generator: %w", err)
			}
			return srv.ServeStdio().Listen(context.Background(), os.Stdin, os.Stdout)
		},
	}
	res.Flags().StringVar(&logFile, "log-file", "mcp.log", "path to log file")
	return res
}
