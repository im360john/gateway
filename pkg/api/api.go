package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/centralmind/gateway/pkg/swaggerator"
	"github.com/doublecloud/transfer/library/go/slices"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/providers"
	"github.com/flowchartsman/swaggerui"
	"net/http"
	"strings"

	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/getkin/kin-openapi/openapi3"
)

// API handles OpenAPI schema generation and sample data serving.
type API struct {
	Schema     abstract.TableMap
	SampleData map[abstract.TableID][]abstract.ChangeItem
	SnapshotF  providers.Snapshot

	transfer *model.Transfer
}

// NewAPI initializes a new API instance.
func NewAPI(
	schema abstract.TableMap,
	sampleData map[abstract.TableID][]abstract.ChangeItem,
	snapshot providers.Snapshot,
	transfer *model.Transfer,
) *API {
	return &API{
		Schema:     schema,
		SampleData: sampleData,
		SnapshotF:  snapshot,

		transfer: transfer,
	}
}

func (api *API) FetchData(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 1 {
		http.Error(w, `{"error": "invalid URL format, expected /{tableName}"}`, http.StatusBadRequest)
		return
	}
	tableName := pathParts[0]

	tid, err := abstract.ParseTableID(tableName)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to parse table name: %s"}`, tableName), http.StatusNotFound)
		return
	}

	storage, err := api.SnapshotF.Storage()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to create storage: %s: %s"}`, tableName, err.Error()), http.StatusBadRequest)
		return
	}
	filter := abstract.WhereStatement("1 = 1")
	for col, value := range r.URL.Query() {
		filter = abstract.WhereStatement(fmt.Sprintf("%s and %s = '%s'", filter, col, value[0]))
	}

	var res abstract.ChangeItem
	if err := storage.LoadTable(context.Background(), abstract.TableDescription{
		Name:   tid.Name,
		Schema: tid.Namespace,
		Filter: filter,
		EtaRow: 0,
		Offset: 0,
	}, func(items []abstract.ChangeItem) error {
		for _, row := range items {
			if !row.IsRowEvent() {
				continue
			}
			res = row
			return nil
		}
		return nil
	}); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to fetch data from storage: %s: %s"}`, tableName, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res.AsMap())
}

func (api *API) SearchData(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 || pathParts[0] != "search" {
		http.Error(w, `{"error": "invalid URL format, expected /search/{tableName}"}`, http.StatusBadRequest)
		return
	}
	tableName := pathParts[1]

	tid, err := abstract.ParseTableID(tableName)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to parse table name: %s"}`, tableName), http.StatusNotFound)
		return
	}

	storage, err := api.SnapshotF.Storage()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to create storage: %s: %s"}`, tableName, err.Error()), http.StatusBadRequest)
		return
	}
	filter := abstract.WhereStatement("1 = 1")
	for col, value := range r.URL.Query() {
		filter = abstract.WhereStatement(fmt.Sprintf("%s and %s = '%s'", filter, col, value[0]))
	}

	var res abstract.ChangeItem
	if err := storage.LoadTable(context.Background(), abstract.TableDescription{
		Name:   tid.Name,
		Schema: tid.Namespace,
		Filter: filter,
		EtaRow: 0,
		Offset: 0,
	}, func(items []abstract.ChangeItem) error {
		for _, row := range items {
			if !row.IsRowEvent() {
				continue
			}
			res = row
			return nil
		}
		return nil
	}); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "unable to fetch data from storage: %s: %s"}`, tableName, err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res.AsMap())
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
	mux.HandleFunc("/sample/", api.ServeSampleData) // Handles /sample/{tableName}
	mux.HandleFunc("/search/", api.SearchData)      // Handles /sample/{tableName}
	mux.HandleFunc("/", api.FetchData)              // Handles /{tableName}
	swagger := swaggerator.Schema(api.Schema)
	raw, _ := json.Marshal(swagger)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(raw)))
}
