package cli

import (
	"github.com/centralmind/gateway/pkg/api"
	"github.com/centralmind/gateway/pkg/logger"
	"github.com/doublecloud/transfer/library/go/core/metrics/solomon"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract/coordinator"
	"github.com/doublecloud/transfer/pkg/providers"
	"github.com/spf13/cobra"
	"net/http"
)

func REST(configPath *string, addr *string) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "REST gateway",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		RunE: func(cmd *cobra.Command, args []string) error {
			transfer, res, err := PrepareConfig(configPath)
			if err != nil {
				return xerrors.Errorf("unable to prepare config: %w", err)
			}
			mux := http.NewServeMux()
			sF, ok := providers.Source[providers.Snapshot](
				logger.NewConsoleLogger(),
				solomon.NewRegistry(solomon.NewRegistryOpts()),
				coordinator.NewFakeClient(),
				transfer,
			)
			if !ok {
				return xerrors.Errorf("no snapshot provider: %T", transfer.Src)
			}
			api.NewAPI(res.Schema, res.Preview, sF, transfer).RegisterRoutes(mux)

			return http.ListenAndServe(*addr, mux)
		},
	}
}
