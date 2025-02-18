package swaggerator

import (
	"github.com/centralmind/gateway/model"
	"github.com/getkin/kin-openapi/openapi3"
)

// Schema dynamically generates an OpenAPI schema based on the given table schema.
func Schema(schema model.Config) *openapi3.T {
	swagger := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       schema.API.Name,
			Description: "Config that dynamically generates accessor for data",
			Version:     schema.API.Version,
		},
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{},
		},
	}

	var paths []openapi3.NewPathsOption
	// Iterate through tables and generate OpenAPI schemas
	for _, info := range schema.Database.Tables {
		schemaProps := make(map[string]*openapi3.SchemaRef)
		for _, col := range info.Columns {
			colType := col.Type

			schemaProps[col.Name] = &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{colType},
				},
			}
		}

		// Add schema to OpenAPI Components
		swagger.Components.Schemas[info.Name] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:       &openapi3.Types{"object"},
				Properties: schemaProps,
			},
		}

		for _, endpoint := range info.Endpoints {
			var endpointParams openapi3.Parameters
			for _, param := range endpoint.Params {
				if param.Location == "" {
					param.Location = "query"
				}
				endpointParams = append(endpointParams, &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     param.Name,
						In:       param.Location,
						Required: param.Required,
						Example:  param.Default,
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:    &openapi3.Types{param.Type},
								Format:  param.Format,
								Default: param.Default,
							},
						},
					},
				})
			}
			paths = append(paths,
				openapi3.WithPath(endpoint.HTTPPath, &openapi3.PathItem{
					Get: &openapi3.Operation{
						Summary:     endpoint.Summary,
						Description: endpoint.Description,
						Tags:        []string{info.Name},
						Parameters:  endpointParams,
						Responses: openapi3.NewResponses(openapi3.WithStatus(200, &openapi3.ResponseRef{
							Value: &openapi3.Response{
								//Description: "JSON object for " + tableName,
								Content: openapi3.Content{
									"application/json": &openapi3.MediaType{
										Schema: &openapi3.SchemaRef{
											Ref: "#/components/schemas/" + info.Name,
										},
									},
								},
							},
						})),
					},
				}),
			)
		}

	}

	swagger.Paths = openapi3.NewPaths(paths...)
	return swagger
}
