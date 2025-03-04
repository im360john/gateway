package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/centralmind/gateway/mcpgenerator"
	gw_model "github.com/centralmind/gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func MCP(configPath *string, addr *string) *cobra.Command {
	cmd := &cobra.Command{
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

			servers, _ := cmd.Flags().GetString("servers")
			serverAddresses := []string{}

			// Add additional servers from the --servers flag if provided
			if servers != "" {
				additionalServers := strings.Split(servers, ",")
				for _, server := range additionalServers {
					serverAddresses = append(serverAddresses, strings.TrimSpace(server))
				}
			}

			if len(serverAddresses) == 0 {
				serverAddresses = append(serverAddresses, fmt.Sprintf("http://localhost%s", *addr))
			}

			srv, err := mcpgenerator.New(*gw)
			if err != nil {
				return xerrors.Errorf("unable to init mcp generator: %w", err)
			}

			logrus.Infof("MCP server is running at: %s/sse", serverAddresses[0])
			return srv.ServeSSE(serverAddresses[0]).Start(*addr)
		},
	}

	cmd.Flags().String("servers", "", "comma-separated list of server addresses")

	return cmd
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
