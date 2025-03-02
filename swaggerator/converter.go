package swaggerator

import (
	"context"
	"embed"
	"github.com/centralmind/gateway/connectors"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/plugins"
	"github.com/danielgtaylor/huma/v2"
)

//go:embed dist
var swagfs embed.FS

// Schema dynamically generates an OpenAPI 3.1 schema based on the given table schema.
func Schema(schema model.Config, addresses ...string) (*huma.OpenAPI, error) {
	api := huma.DefaultConfig(schema.API.Name, "3.1.0").OpenAPI
	api.Info.Title = schema.API.Name
	api.Info.Description = "Config that dynamically generates accessor for data"
	api.Info.Version = schema.API.Version
	api.Paths = make(map[string]*huma.PathItem)

	connector, err := connectors.New(schema.Database.Type, schema.Database.Connection)
	if err != nil {
		return nil, xerrors.Errorf("unable to init connector: %w", err)
	}

	// Add all server addresses
	for i, address := range addresses {
		var description string
		switch {
		case i == 0 && strings.HasPrefix(address, "http://localhost"):
			description = "Local development server"
		case strings.Contains(address, "dev"):
			description = "Development server"
		case strings.Contains(address, "stage"):
			description = "Staging server"
		case strings.Contains(address, "prod"):
			description = "Production server"
		default:
			description = "Server " + address
		}
		api.Servers = append(api.Servers, &huma.Server{
			URL:         address,
			Description: description,
		})
	}

	// If no servers were provided, add a default one
	if len(api.Servers) == 0 {
		api.Servers = append(api.Servers, &huma.Server{
			URL:         "http://localhost:9090",
			Description: "localhost",
		})
	}

	// Iterate through tables and generate OpenAPI schemas
	for _, info := range schema.Database.Tables {
		for _, endpoint := range info.Endpoints {
			cols, err := connector.InferQuery(context.Background(), endpoint.Query)
			if err != nil {
				logrus.Warnf("unable to infer query %s: %v", endpoint.Query, err)
			}
			schemaProps := map[string]*huma.Schema{}
			for _, col := range cols {
				schemaProps[col.Name] = &huma.Schema{
					Type: string(col.Type),
				}
				if col.Type == model.TypeDatetime {
					schemaProps[col.Name] = &huma.Schema{
						Type:   "string",
						Format: "date-time",
					}
				}
			}

			var params []*huma.Param
			for _, param := range endpoint.Params {
				if param.Location == "" {
					param.Location = "query"
				}
				params = append(params, &huma.Param{
					Name:     param.Name,
					In:       param.Location,
					Required: param.Required,
					Schema: &huma.Schema{
						Type:    param.Type,
						Format:  param.Format,
						Default: param.Default,
					},
				})
			}
			resSchema := &huma.Schema{
				Type:       "object",
				Properties: schemaProps,
			}
			if endpoint.IsArrayResult {
				resSchema = &huma.Schema{
					Type:  "array",
					Items: resSchema,
				}
			}
			operation := &huma.Operation{
				Summary:     endpoint.Summary,
				Description: endpoint.Description,
				OperationID: endpoint.MCPMethod,
				Tags:        []string{info.Name},
				Parameters:  params,
				Responses: map[string]*huma.Response{
					"200": {
						Description: "Success",
						Content: map[string]*huma.MediaType{
							"application/json": {
								Schema: resSchema,
							},
						},
					},
					"404": {
						Description: "Not Found",
						Content: map[string]*huma.MediaType{
							"application/json": {
								Schema: &huma.Schema{},
							},
						},
					},
					"500": {
						Description: "Error",
						Content: map[string]*huma.MediaType{
							"application/json": {
								Schema: &huma.Schema{},
							},
						},
					},
				},
			}
			api.Paths[endpoint.HTTPPath] = &huma.PathItem{}
			switch endpoint.HTTPMethod {
			case "GET":
				api.Paths[endpoint.HTTPPath].Get = operation
			case "DELETE":
				api.Paths[endpoint.HTTPPath].Delete = operation
			case "POST":
				api.Paths[endpoint.HTTPPath].Post = operation
			case "PATCH":
				api.Paths[endpoint.HTTPPath].Patch = operation
			case "PUT":
				api.Paths[endpoint.HTTPPath].Put = operation
			}
		}
	}

	api, err = plugins.Enrich(schema.Plugins, api)
	if err != nil {
		return nil, xerrors.Errorf("unable to enrich swagger schema: %w", err)
	}
	return api, nil
}

func byteHandler(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Write(b)
	}
}

// Handler returns a handler that will serve a self-hosted Swagger UI with your spec embedded
func Handler(spec []byte) http.Handler {
	// render the index template with the proper spec name inserted
	static, _ := fs.Sub(swagfs, "dist")
	mux := http.NewServeMux()
	mux.HandleFunc("/swagger_spec", byteHandler(spec))
	mux.Handle("/", http.FileServer(http.FS(static)))
	return mux
}
