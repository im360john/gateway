package rawmcp

import (
	"context"
	"fmt"
	"github.com/centralmind/gateway/prompter"
	"golang.org/x/xerrors"
	"strings"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/mcp"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/plugins"
	"github.com/centralmind/gateway/server"
)

type Server struct {
	server *server.MCPServer
}

func New(
	schema model.Config,
) (*Server, error) {
	srv := server.NewMCPServer(
		"mcp-data-gateway",
		"0.0.1",
		server.WithInstructions(fmt.Sprintf(`Server for data gateway

Tools actions flow if follows:
0. (optional) - take a look on tables to check what data shall be discovered
1. discover_data, this will return a description of data available along side with data samples
2. prepare_query, this will return wheather a query is working or not, if not working - need to fix according to error mesasge
3. query, this will return an actual data set for a SQL query

!Important rules:
	- SQL queries must be Pure SQL that will be used in golang SQLx on top of database - %[1]s and be fully parameterized (using named parameters) to prevent SQL injection.
`, schema.Database.Type)),
	)

	var interceptors []plugins.Interceptor
	for k, v := range schema.Plugins {
		plugin, err := plugins.New(k, v)
		if err != nil {
			return nil, err
		}
		interceptor, ok := plugin.(plugins.Interceptor)
		if !ok {
			continue
		}
		interceptors = append(interceptors, interceptor)
	}
	connector, err := connectors.New(schema.Database.Type, schema.Database.Connection)
	if err != nil {
		return nil, xerrors.Errorf("unable to init connector: %w", err)
	}
	connector, err = plugins.Wrap(schema.Plugins, connector)
	if err != nil {
		return nil, xerrors.Errorf("unable to init connector plugins: %w", err)
	}

	srv.AddTool(mcp.NewTool(
		"list_tables",
		mcp.WithDescription(`Return list of tables that available for data.
This is usually first this agent shall call.
`),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, err := connector.Discovery(ctx)
		if err != nil {
			return nil, xerrors.Errorf("unable to discover data: %w", err)
		}

		var content []mcp.Content
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Found %v records-(s).", len(data)),
		})
		for _, record := range data {
			schema := prompter.SchemaFromConfig(connector.Config())
			if schema != "" {
				record.Name = fmt.Sprintf("%v.%v", schema, record.Name)
			}
			content = append(content, mcp.TextContent{
				Type: "text",
				Text: prompter.Yamlify(record),
			})
		}
		return &mcp.CallToolResult{
			Content: content,
		}, nil
	})
	srv.AddTool(mcp.NewTool(
		"discover_data",
		mcp.WithDescription(`Discover data structure for connected gateway.
tables_list parameter is comma separated table to fetch data samples.
Disovery better to call with a list of interested tables, since it will load all their samples.
`),
		mcp.WithString("tables_list"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, err := connector.Discovery(ctx)
		if err != nil {
			return nil, xerrors.Errorf("unable to discover data: %w", err)
		}
		var content []mcp.Content
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Found a %v tables-(s).", len(data)),
		})
		allTables, err := connector.Discovery(ctx)
		if err != nil {
			return nil, xerrors.Errorf("unable to discover all tables: %w", err)
		}
		tablesList := request.Params.Arguments["tables_list"].(string)

		tableSet := map[string]bool{}

		for _, table := range strings.Split(tablesList, ",") {
			tableSet[table] = true
		}
		if len(tablesList) == 0 {
			for _, table := range allTables {
				tableSet[table.Name] = true
			}
		}
		var tablesToGenerate []prompter.TableData
		for _, table := range allTables {
			if !tableSet[table.Name] {
				continue
			}
			sample, err := connector.Sample(ctx, table)
			if err != nil {
				return nil, xerrors.Errorf("unable to discover sample: %w", err)
			}
			tablesToGenerate = append(tablesToGenerate, prompter.TableData{
				Columns:  table.Columns,
				Name:     table.Name,
				Sample:   sample,
				RowCount: table.RowCount,
			})
		}

		content = append(content, mcp.TextContent{
			Type: "text",
			Text: prompter.TablesPrompt(tablesToGenerate, prompter.SchemaFromConfig(connector.Config())),
		})

		return &mcp.CallToolResult{
			Content: content,
		}, nil
	})
	srv.AddTool(mcp.NewTool(
		"prepare_query",
		mcp.WithDescription(`Verify query and prepare output structure for query.
This tool shall be executed before query, to examine output structure and verify that query is correct.
`),
		mcp.WithString("query", mcp.Required()),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resSchema, err := connector.InferQuery(ctx, request.Params.Arguments["query"].(string))
		if err != nil {
			return nil, xerrors.Errorf("unable to infer query: %w", err)
		}
		var content []mcp.Content
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Query has a %v column-(s).", len(resSchema)),
		})
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: prompter.Yamlify(resSchema),
		})
		return &mcp.CallToolResult{
			Content: content,
		}, nil
	})
	srv.AddTool(mcp.NewTool(
		"query",
		mcp.WithDescription("Query data structure for connected gateway"),
		mcp.WithString("query", mcp.Required()),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resData, err := connector.Query(
			ctx,
			model.Endpoint{Query: request.Params.Arguments["query"].(string)},
			make(map[string]any),
		)
		if err != nil {
			return nil, xerrors.Errorf("unable to infer query: %w", err)
		}
		var content []mcp.Content
		content = append(content, mcp.TextContent{
			Type: "text",
			Text: fmt.Sprintf("Found %v records-(s).", len(resData)),
		})
		for _, record := range resData {
			content = append(content, mcp.TextContent{
				Type: "text",
				Text: prompter.Yamlify(record),
			})
		}
		return &mcp.CallToolResult{
			Content: content,
		}, nil
	})

	return &Server{
		server: srv,
	}, nil
}

func (s *Server) ServeSSE(addr string) *server.SSEServer {
	return server.NewSSEServer(s.server, addr)
}

func (s *Server) ServeStdio() *server.StdioServer {
	return server.NewStdioServer(s.server)
}
