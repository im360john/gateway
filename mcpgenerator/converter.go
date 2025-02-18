package mcpgenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/plugins"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server *server.MCPServer
}

func New(
	schema model.Config,
) (*MCPServer, error) {
	srv := server.NewMCPServer("mcp-data-gateway", "0.0.1")
	var interceptors []plugins.Interceptor
	for k, v := range schema.Plugins {
		interceptor, err := plugins.New(k, v)
		if err != nil {
			return nil, err
		}
		interceptors = append(interceptors, interceptor)
	}
	for _, info := range schema.Database.Tables {
		for _, endpoint := range info.Endpoints {
			var opts []mcp.ToolOption
			for _, col := range endpoint.Params {
				if col.Required {
					opts = append(opts, ArgumentOption(col, mcp.Required()))
				} else {
					opts = append(opts, ArgumentOption(col))
				}
			}

			srv.AddTool(mcp.NewTool(
				endpoint.MCPMethod,
				opts...,
			), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

				return nil, nil
			})
		}
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

func ArgumentOption(col model.EndpointParams, opts ...mcp.PropertyOption) mcp.ToolOption {
	opts = append(opts, mcp.Title(fmt.Sprintf("Column %s", col.Name)))
	opts = append(opts, func(m map[string]interface{}) {
		m["default"] = col.Default
	})

	switch col.Type {
	case "integer", "double", "float", "number":
		return mcp.WithNumber(col.Name, opts...)

	case "string":
		return mcp.WithString(col.Name, opts...)

	case "bool", "boolean":
		return mcp.WithBoolean(col.Name, opts...)

	default:
		return mcp.WithString(col.Name)
	}
}
