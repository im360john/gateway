package castx

import (
	"github.com/centralmind/gateway/model"
	"github.com/spf13/cast"
	"golang.org/x/xerrors"
)

func ParamsE(endpoint model.Endpoint, params map[string]any) (map[string]any, error) {
	processedParams := make(map[string]any)
	for _, param := range endpoint.Params {
		if _, ok := params[param.Name]; !ok {
			continue
		}
		switch param.Type {
		case "string":
			processedParams[param.Name] = cast.ToString(params[param.Name])
		case "number":
			var err error
			processedParams[param.Name], err = cast.ToIntE(params[param.Name])
			if err != nil {
				processedParams[param.Name], err = cast.ToFloat64E(params[param.Name])
				if err != nil {
					return nil, xerrors.Errorf("unable to parse number: %s: %w", param.Name, err)
				}
			}
		case "bool", "boolean":
			processedParams[param.Name] = cast.ToBool(params[param.Name])
		default:
			processedParams[param.Name] = cast.ToString(params[param.Name])
		}
	}
	return processedParams, nil
}

func Process(row map[string]any) map[string]any {
	for k := range row {
		if bb, ok := row[k].([]byte); ok {
			row[k] = string(bb)
		}
	}
	return row
}
