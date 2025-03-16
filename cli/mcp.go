package cli

import (
	"context"
	"github.com/centralmind/gateway/logger"
	"github.com/centralmind/gateway/mcpgenerator"
	gw_model "github.com/centralmind/gateway/model"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
)

func Stdio(configPath *string) *cobra.Command {
	var logFile string
	var rawMode bool
	res := &cobra.Command{
		Use:   "stdio",
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
			if rawMode {
				if err := mcpgenerator.AddRawProtocol(*gw, srv.Server()); err != nil {
					return xerrors.Errorf("unable to add raw mcp protocol: %w", err)
				}
			}

			return srv.ServeStdio().Listen(context.Background(), os.Stdin, os.Stdout)
		},
	}

	res.Flags().BoolVar(&rawMode, "raw", false, "enable as raw protocol")
	res.Flags().StringVar(&logFile, "log-file", filepath.Join(logger.DefaultLogDir(), "mcp.log"), "path to log file")
	return res
}
