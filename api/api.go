package api

import (
	"context"
	"encoding/json"
	"github.com/centralmind/gateway/connectors"
	gw_model "github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/plugins"
	"github.com/centralmind/gateway/swaggerator"
	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
)

// API handles OpenAPI schema generation and sample data serving.
type API struct {
	Schema       gw_model.Config
	interceptors []plugins.Interceptor
	connector    connectors.Connector
}

// NewAPI initializes a new API instance.
func NewAPI(
	schema gw_model.Config,
) (*API, error) {
	var interceptors []plugins.Interceptor
	for k, v := range schema.Plugins {
		interceptor, err := plugins.New(k, v)
		if err != nil {
			return nil, err
		}
		interceptors = append(interceptors, interceptor)
	}
	connector, err := connectors.New(schema.Gateway.Type, schema.Gateway.Connection)
	if err != nil {
		return nil, errors.Errorf("unable to init connector: %w", err)
	}
	if err := connector.Ping(context.Background()); err != nil {
		return nil, errors.Errorf("unable to ping: %w", err)
	}
	return &API{
		Schema:       schema,
		interceptors: interceptors,
		connector:    connector,
	}, nil
}

// RegisterRoutes registers API endpoints.
func (api *API) RegisterRoutes(mux *http.ServeMux) {
	swagger := swaggerator.Schema(api.Schema)
	raw, _ := json.Marshal(swagger)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(raw)))
	r := gin.Default()
	for _, table := range api.Schema.Gateway.Tables {
		for _, endpoint := range table.Endpoints {
			r.Handle(endpoint.HTTPMethod, convertSwaggerToGin(endpoint.HTTPPath), api.Handler(endpoint))
		}
	}
	mux.Handle("/", r.Handler())
}

func (api *API) Handler(endpoint gw_model.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := make(map[string]any)

		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		for key, values := range c.Request.URL.Query() {
			if len(values) == 1 {
				params[key] = values[0]
			} else {
				params[key] = values
			}
		}
		for _, param := range endpoint.Params {
			if _, ok := params[param.Name]; !ok {
				params[param.Name] = nil
			}
		}

		raw, err := api.connector.Query(c.Request.Context(), endpoint, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var res []map[string]any
	MAIN:
		for _, row := range raw {
			for _, interceptor := range api.interceptors {
				r, skip := interceptor.Process(row, c.Request.Header)
				if skip {
					continue MAIN
				}
				row = r
			}
			res = append(res, row)
		}
		c.JSON(http.StatusOK, res)
	}
}

func convertSwaggerToGin(swaggerURL string) string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	return re.ReplaceAllString(swaggerURL, ":$1")
}
