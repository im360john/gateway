package cli

import (
	"context"
	gw_model "github.com/centralmind/gateway/pkg/model"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/worker/tasks"
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

	endpoint, err := gw.Endpoint()
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("endpoint is invalid: %w", err)
	}
	transfer := &model.Transfer{
		Type:           abstract.TransferTypeSnapshotOnly,
		Src:            endpoint,
		Dst:            new(model.MockDestination),
		DataObjects:    nil,
		Transformation: nil,
	}
	res := tasks.TestEndpoint(context.Background(), &tasks.TestEndpointParams{
		Transfer:             transfer,
		TransformationConfig: nil,
		ParamsSrc: &tasks.EndpointParam{
			Type:  transfer.SrcType(),
			Param: transfer.SrcJSON(),
		},
		ParamsDst: nil,
	}, abstract.NewTestResult())
	if res.Err() != nil {
		return nil, nil, nil, xerrors.Errorf("unable to runRest source: %w", res.Err())
	}
	return transfer, gw, res, nil
}
