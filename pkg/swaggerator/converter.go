package swaggerator

import (
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/getkin/kin-openapi/openapi3"
	"go.ytsaurus.tech/yt/go/schema"
	"strings"
)

// Schema dynamically generates an OpenAPI schema based on the given table schema.
func Schema(schema abstract.TableMap) *openapi3.T {
	swagger := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Data Gateway Schema API",
			Description: "Gateway that dynamically generates accessor for data",
			Version:     "1.0.0",
		},
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{},
		},
	}

	// Iterate through tables and generate OpenAPI schemas
	for tid, info := range schema {
		schemaProps := make(map[string]*openapi3.SchemaRef)
		var queryParams openapi3.Parameters
		tableName := strings.ReplaceAll(tid.Fqtn(), "\"", "")

		for _, col := range info.Schema.Columns() {
			colType := mapDataTypeToOpenAPI(col.DataType)

			schemaProps[col.ColumnName] = &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: &openapi3.Types{colType},
				},
			}

			if col.IsKey() {
				queryParams = append(queryParams, &openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:     col.ColumnName,
						In:       "query",
						Required: false,
						Schema: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: &openapi3.Types{mapDataTypeToOpenAPI(col.DataType)},
							},
						},
					},
				})
			}
		}

		// Add schema to OpenAPI Components
		swagger.Components.Schemas[tableName] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:       &openapi3.Types{"object"},
				Properties: schemaProps,
			},
		}

		swagger.Paths = openapi3.NewPaths(
			openapi3.WithPath("/sample/"+tableName, &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "Get JSON example for " + tableName,
					Description: "Generates example JSON based on " + tableName,
					Tags:        []string{tableName},
					Responses: openapi3.NewResponses(openapi3.WithStatus(200, &openapi3.ResponseRef{
						Value: &openapi3.Response{
							//Description: "JSON object for " + tableName,
							Content: openapi3.Content{
								"application/json": &openapi3.MediaType{
									Schema: &openapi3.SchemaRef{
										Ref: "#/components/schemas/" + tableName,
									},
								},
							},
						},
					})),
				},
			}),
			openapi3.WithPath("/"+tableName, &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "Get one row for " + tableName,
					Description: "Get row in json format based on " + tableName,
					Tags:        []string{tableName},
					Parameters:  queryParams,
					Responses: openapi3.NewResponses(openapi3.WithStatus(200, &openapi3.ResponseRef{
						Value: &openapi3.Response{
							//Description: "JSON object for " + tableName,
							Content: openapi3.Content{
								"application/json": &openapi3.MediaType{
									Schema: &openapi3.SchemaRef{
										Ref: "#/components/schemas/" + tableName,
									},
								},
							},
						},
					})),
				},
			}),
		)
	}

	return swagger
}

// mapDataTypeToOpenAPI converts column data types to OpenAPI types.
func mapDataTypeToOpenAPI(dataType string) string {
	switch dataType {
	case schema.TypeInt64.String(), schema.TypeInt32.String(), schema.TypeInt16.String(), schema.TypeInt8.String(),
		schema.TypeUint64.String(), schema.TypeUint32.String(), schema.TypeUint16.String(), schema.TypeUint8.String():
		return "integer"

	case schema.TypeFloat32.String(), schema.TypeFloat64.String():
		return "number"

	case schema.TypeBytes.String(), schema.TypeString.String():
		return "string"

	case schema.TypeBoolean.String():
		return "boolean"

	case schema.TypeAny.String():
		return "object"

	case schema.TypeDate.String(), schema.TypeDatetime.String(), schema.TypeTimestamp.String(), schema.TypeInterval.String():
		return "string" // Dates are usually represented as ISO 8601 strings

	default:
		return "string"
	}
}
