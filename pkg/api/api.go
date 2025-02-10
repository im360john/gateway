package api

import (
	"encoding/json"
	"fmt"
	"github.com/centralmind/gateway/pkg/swaggerator"
	"github.com/doublecloud/transfer/library/go/slices"
	"net/http"
	"strings"

	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/getkin/kin-openapi/openapi3"
)

// API handles OpenAPI schema generation and sample data serving.
type API struct {
	Schema     abstract.TableMap
	SampleData map[abstract.TableID][]abstract.ChangeItem
}

// NewAPI initializes a new API instance.
func NewAPI(schema abstract.TableMap, sampleData map[abstract.TableID][]abstract.ChangeItem) *API {
	return &API{
		Schema:     schema,
		SampleData: sampleData,
	}
}

// ServeSwaggerJSON serves the generated OpenAPI JSON.
func (api *API) ServeSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	swagger := swaggerator.Schema(api.Schema)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(swagger)
}

// ServeSampleData serves example data for a given table.
func (api *API) ServeSampleData(w http.ResponseWriter, r *http.Request) {
	// Extract table name from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "sample" {
		http.Error(w, `{"error": "invalid URL format, expected /sample/{tableName}"}`, http.StatusBadRequest)
		return
	}
	tableName := pathParts[1]

	tid, err := abstract.ParseTableID(tableName)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to parse table name: %s"}`, tableName), http.StatusNotFound)
		return
	}

	sample, exists := api.SampleData[*tid]
	if !exists {
		http.Error(w, `{"error": "sample data not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(slices.Map(sample, func(t abstract.ChangeItem) any {
		return t.AsMap()
	}))
}

// generateSwaggerJSON generates OpenAPI JSON (stub implementation).
func (api *API) generateSwaggerJSON() *openapi3.T {
	// TODO: Implement dynamic OpenAPI generation using api.Schema
	return nil
}

// RegisterRoutes registers API endpoints.
func (api *API) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swagger.json", api.ServeSwaggerJSON)
	mux.HandleFunc("/sample/", api.ServeSampleData) // Handles /sample/{tableName}
}
