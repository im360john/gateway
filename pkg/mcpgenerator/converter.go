package mcpgenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/providers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.ytsaurus.tech/yt/go/schema"
	"strings"
)

type MCPServer struct {
	server *server.MCPServer
}

func New(
	schema abstract.TableMap,
	sampleData map[abstract.TableID][]abstract.ChangeItem,
	snapshot providers.Snapshot,
	transfer *model.Transfer,
) (*MCPServer, error) {
	srv := server.NewMCPServer("mcp-data-gateway", "0.0.1")

	for tid, info := range schema {
		tableName := strings.ReplaceAll(tid.Fqtn(), "\"", "")

		var requestOpts []mcp.ToolOption
		requestOpts = append(requestOpts, mcp.WithDescription(fmt.Sprintf("Get exact record from %s", tableName)))
		for _, col := range info.Schema.Columns() {
			if col.IsKey() {
				requestOpts = append(requestOpts, ColToolOption(col))
			}
		}
		srv.AddTool(mcp.NewTool(
			"get_"+tableName,
			requestOpts...,
		), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			arguments := request.Params.Arguments
			filter := abstract.WhereStatement("1 = 1")
			for col, value := range arguments {
				filter = abstract.WhereStatement(fmt.Sprintf("%s and %s = '%v'", filter, col, value))
			}
			storage, err := snapshot.Storage()
			if err != nil {
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Unable to construct storage: %w", err),
						},
					},
					IsError: true,
				}, nil
			}

			var found bool
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
					found = true
					return nil
				}
				return nil
			}); err != nil {
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Unable to find row: %w", err),
						},
					},
					IsError: true,
				}, nil
			}
			if !found {
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Row: %v not found", request.Params.Arguments),
						},
					},
					IsError: false,
				}, nil
			}
			return &mcp.CallToolResult{
				Content: []interface{}{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("Found a row: %v in %s.", request.Params.Arguments, tableName),
					},
					mcp.TextContent{
						Annotated: mcp.Annotated{},
						Type:      "test",
						Text:      jsonify(res.AsMap()),
					},
				},
			}, nil
		})
	}

	return &MCPServer{
		server: srv,
	}, nil
}

func jsonify(asMap map[string]interface{}) string {
	res, _ := json.MarshalIndent(asMap, "", "  ")
	return string(res)
}

func (s *MCPServer) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server, fmt.Sprintf("http://%s", addr))
}

func (s *MCPServer) ServeStdio() *server.StdioServer {
	return server.NewStdioServer(s.server)
}

func ColToolOption(col abstract.ColSchema) mcp.ToolOption {
	var opts []mcp.PropertyOption
	if col.Required {
		opts = append(opts, mcp.Required())
	}
	opts = append(opts, mcp.Title(fmt.Sprintf("Column %s", col.ColumnName)))

	switch col.DataType {
	case schema.TypeInt64.String(), schema.TypeInt32.String(), schema.TypeInt16.String(), schema.TypeInt8.String(),
		schema.TypeUint64.String(), schema.TypeUint32.String(), schema.TypeUint16.String(), schema.TypeUint8.String():
		return mcp.WithNumber(col.ColumnName, opts...)

	case schema.TypeFloat32.String(), schema.TypeFloat64.String():
		return mcp.WithNumber(col.ColumnName, opts...)

	case schema.TypeBytes.String(), schema.TypeString.String():
		return mcp.WithString(col.ColumnName, opts...)

	case schema.TypeBoolean.String():
		return mcp.WithBoolean(col.ColumnName, opts...)

	case schema.TypeAny.String():
		return mcp.WithString(col.ColumnName, opts...)

	case schema.TypeDate.String(), schema.TypeDatetime.String(), schema.TypeTimestamp.String(), schema.TypeInterval.String():
		return mcp.WithString(col.ColumnName, opts...) // Dates are usually represented as ISO 8601 strings

	default:
		return mcp.WithString(col.ColumnName)
	}
}
