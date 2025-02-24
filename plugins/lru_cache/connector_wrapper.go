package lrucache

import (
	"context"
	"fmt"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/model"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"sort"
	"strings"
)

type Connector struct {
	connectors.Connector

	config Config
	lru    *expirable.LRU[string, []map[string]any]
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	key := keyify(endpoint, params)
	v, ok := c.lru.Get(key)
	if !ok {
		res, err := c.Connector.Query(ctx, endpoint, params)
		if err != nil {
			return nil, err
		}
		_ = c.lru.Add(key, res)
		return res, nil
	}
	return v, nil
}

func keyify(endpoint model.Endpoint, params map[string]any) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sb := strings.Builder{}
	for _, k := range keys {
		_, _ = sb.WriteString(fmt.Sprintf("%v", params[k]))
	}
	return endpoint.MCPMethod + "/" + sb.String()
}
