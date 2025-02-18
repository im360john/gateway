package mcpgenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doublecloud/transfer/pkg/abstract"
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
	snapshot providers.Snapshot,
) (*MCPServer, error) {
	srv := server.NewMCPServer("mcp-data-gateway", "0.0.1")

	for tid, info := range schema {
		tableName := strings.ReplaceAll(
			strings.ReplaceAll(
				tid.Fqtn(),
				"\"",
				"",
			),
			".",
			"_",
		)
		tableName = strings.TrimLeft(tableName, "public_")

		var allOptions []mcp.ToolOption
		allOptions = append(allOptions, mcp.WithDescription(fmt.Sprintf("Find record from %s table, all columns is queriable", tableName)))
		var keyOptions []mcp.ToolOption
		keyOptions = append(keyOptions, mcp.WithDescription(fmt.Sprintf("Get exact record from %s", tableName)))
		for _, col := range info.Schema.Columns() {
			if col.IsKey() {
				keyOptions = append(keyOptions, ArgumentOption(col, mcp.Required()))
			}
			allOptions = append(allOptions, ArgumentOption(col), OperandOption(col))
		}
		var findRecord = func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			arguments := request.Params.Arguments
			filter := abstract.WhereStatement("1 = 1")
			for col, value := range arguments {
				if strings.HasSuffix(col, "_operand") {
					continue
				}
				operand := " = "
				if reqOp, ok := arguments[fmt.Sprintf("%s_operand", col)]; ok {
					switch reqOp {
					case "more":
						reqOp = " > "
					case "less":
						reqOp = " < "
					}
				}
				filter = abstract.WhereStatement(fmt.Sprintf("%s and %s %s '%v'", filter, col, operand, value))
			}
			storage, err := snapshot.Storage()
			if err != nil {
				return &mcp.CallToolResult{
					Content: []interface{}{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("Unable to construct storage: %s", err),
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
							Text: fmt.Sprintf("Unable to find row: %s", err),
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
						Type:      "text",
						Text:      jsonify(res.AsMap()),
					},
				},
			}, nil
		}
		srv.AddTool(mcp.NewTool(
			fmt.Sprintf("find_%s", tableName),
			allOptions...,
		), findRecord)
		srv.AddTool(mcp.NewTool(
			fmt.Sprintf("get_%s", tableName),
			keyOptions...,
		), findRecord)
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

func OperandOption(col abstract.ColSchema) mcp.ToolOption {
	return mcp.WithString(
		fmt.Sprintf("%s_operator", col.ColumnName),
		mcp.Enum("equal", "more", "less"),
		mcp.DefaultString("equal"),
		mcp.Description(fmt.Sprintf("what operator apply to query by `%s` field", col.ColumnName)),
	)
}

func ArgumentOption(col abstract.ColSchema, opts ...mcp.PropertyOption) mcp.ToolOption {
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
