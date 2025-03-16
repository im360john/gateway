package cli

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	gw_model "github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/restgenerator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func REST(configPath *string, addr *string, servers *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest",
		Short: "REST gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) error {
			disableSwagger, _ := cmd.Flags().GetBool("disable-swagger")
			prefix, _ := cmd.Flags().GetString("prefix")

			gwRaw, err := os.ReadFile(*configPath)
			if err != nil {
				return xerrors.Errorf("unable to read yaml config file: %w", err)
			}
			gw, err := gw_model.FromYaml(gwRaw)
			if err != nil {
				return xerrors.Errorf("unable to parse config file: %w", err)
			}
			mux := http.NewServeMux()
			a, err := restgenerator.New(*gw, prefix)

			if err != nil {
				return xerrors.Errorf("unable to init api: %w", err)
			}

			// Create the list of server addresses for RegisterRoutes
			serverAddresses := []string{}

			// Add additional servers from the --servers flag if provided
			if *servers != "" {
				additionalServers := strings.Split(*servers, ",")
				for _, server := range additionalServers {
					serverAddresses = append(serverAddresses, strings.TrimSpace(server))
				}
			}

			if len(serverAddresses) == 0 {
				serverAddresses = append(serverAddresses, fmt.Sprintf("http://localhost%s", *addr))
			}

			// Use the disable-swagger flag value passed from parent command
			// Register routes with all server addresses and disable-swagger flag
			if err := a.RegisterRoutes(mux, disableSwagger, serverAddresses...); err != nil {
				return err
			}

			if !disableSwagger {
				logrus.Infof("Docs available at: %s/%s", serverAddresses[0], prefix)
			}

			return http.ListenAndServe(*addr, mux)
		},
	}

	cmd.Flags().Bool("disable-swagger", false, "disable Swagger UI")
	cmd.Flags().String("prefix", "", "prefix for API endpoints")

	return cmd
}
