package luarls

// Config represents Lua Row-Level Security configuration
type Config struct {
	// Script is the Lua script content for RLS logic
	Script string `yaml:"script"`

	// ScriptFile is the path to Lua script file (alternative to Script)
	ScriptFile string `yaml:"script_file"`

	// Variables defines global variables available to Lua script
	Variables map[string]interface{} `yaml:"variables"`

	// CacheSize is the size of the Lua VM cache
	CacheSize int `yaml:"cache_size"`
}

func (c Config) Tag() string {
	return "lua_rls"
}

func (c Config) Doc() string {
	return docString
}
