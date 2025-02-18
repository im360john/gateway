package cli

import (
	gw_model "github.com/centralmind/gateway/model"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"os"
)

func PrepareConfig(configPath *string) (*model.Transfer, *gw_model.Config, *abstract.TestResult, error) {
	gwRaw, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("unable to read yaml config file: %w", err)
	}
	gw, err := gw_model.FromYaml(gwRaw)
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("unable to parse config file: %w", err)
	}

	return nil, gw, nil, nil
}
