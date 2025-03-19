package mcpgenerator

import (
	"encoding/json"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/plugins"
	"github.com/centralmind/gateway/server"
	"golang.org/x/xerrors"
	"sync"
)

type MCPServer struct {
	server    *server.MCPServer
	connector connectors.Connector
	tools     []model.Endpoint

	mu sync.Mutex
}

func New(
	db model.Database,
	plugs map[string]any,
) (*MCPServer, error) {
	srv := server.NewMCPServer("mcp-data-gateway", "0.0.1")
	var interceptors []plugins.Interceptor
	for k, v := range plugs {
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
	connector, err := connectors.New(db.Type, db.Connection)
	if err != nil {
		return nil, xerrors.Errorf("unable to init connector: %w", err)
	}
	connector, err = plugins.Wrap(plugs, connector)
	if err != nil {
		return nil, xerrors.Errorf("unable to init connector plugins: %w", err)
	}

	return &MCPServer{
		server:    srv,
		connector: connector,
	}, nil
}

func (s *MCPServer) ServeSSE(addr string, prefix string) *server.SSEServer {
	return server.NewSSEServer(s.server, addr, prefix)
}

func (s *MCPServer) ServeStdio() *server.StdioServer {
	return server.NewStdioServer(s.server)
}

func (s *MCPServer) Server() *server.MCPServer {
	return s.server
}

func jsonify(data any) string {
	res, _ := json.Marshal(data)
	return string(res)
}
