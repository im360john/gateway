package api

import (
	"encoding/json"
	"github.com/centralmind/gateway/pkg/logger"
	gw_model "github.com/centralmind/gateway/pkg/model"
	"github.com/centralmind/gateway/pkg/plugins"
	"github.com/centralmind/gateway/pkg/swaggerator"
	"github.com/doublecloud/transfer/pkg/providers/postgres"
	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"
	"regexp"
)

// API handles OpenAPI schema generation and sample data serving.
type API struct {
	Schema       gw_model.Config
	interceptors []plugins.Interceptor
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
	return &API{
		Schema:       schema,
		interceptors: interceptors,
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

		var src postgres.PgSource
		if err := json.Unmarshal([]byte(api.Schema.ParamRaw()), &src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cfg, err := postgres.MakeConnConfigFromStorage(logger.NewConsoleLogger(), src.ToStorageParams(nil))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		db := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
		rows, err := db.NamedQueryContext(c.Request.Context(), endpoint.Query, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		res := make([]map[string]any, 0)
		for rows.Next() {
			row := map[string]any{}
			if err := rows.MapScan(row); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, interceptor := range api.interceptors {
				r, skip := interceptor.Process(row, c.Request.Header)
				if skip {
					continue
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
