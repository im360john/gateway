package luarls

import (
	"fmt"
	"github.com/centralmind/gateway/plugins"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	plugins.Register(func(cfg Config) (plugins.Interceptor, error) {
		return New(cfg)
	})
}

type Plugin struct {
	script string
}

func (l Plugin) Doc() string {
	return `
Allow to execute lua script for every row in result set

# Example YAML configuration:

lua_rls:
  script: |
    function filter_rows(row, context)
      if context.user_role == "admin" then
        return true
      end
      return row.tenant_id == context.tenant_id
    end
  variables:
    max_rows: 1000
    debug: true
  cache_size: 100
`
}

func (l Plugin) Process(row map[string]any, headers map[string][]string) (processed map[string]any, skipped bool) {
	st := lua.NewState()
	if err := st.DoString(l.script); err != nil {
		return nil, false
	}
	fn := st.GetGlobal("check_visibility")
	if fn == lua.LNil {
		return nil, false
	}
	defer st.Close()

	rowTable := st.NewTable()
	for k, v := range row {
		st.SetTable(rowTable, lua.LString(k), lua.LString(fmt.Sprintf("%v", v)))
	}

	// Преобразуем headers в lua-таблицу
	headersTable := st.NewTable()
	for k, v := range headers {
		if len(v) != 1 {
			continue
		}
		st.SetTable(headersTable, lua.LString(k), lua.LString(v[0]))
	}

	if err := st.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, rowTable, headersTable); err != nil {
		return nil, true
	}

	ret := st.Get(-1)
	st.Pop(1)

	return row, lua.LVAsBool(ret)
}

func New(config Config) (plugins.Interceptor, error) {
	st := lua.NewState()
	if err := st.DoString(config.Script); err != nil {
		return nil, err
	}
	fn := st.GetGlobal("check_visibility")
	if fn == lua.LNil {
		return nil, fmt.Errorf("entry point check_visibility not found")
	}
	st.Close()
	return &Plugin{
		script: config.Script,
	}, nil
}
