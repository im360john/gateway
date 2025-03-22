//go:build cgo

package main

import (
	_ "github.com/centralmind/gateway/connectors/duckdb"
	_ "github.com/centralmind/gateway/connectors/sqlite"
)
