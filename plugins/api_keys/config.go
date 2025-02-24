package api_keys

type Config struct {
	Name     string `yaml:"name"`
	Keys     []Key  `yaml:"keys"`
	KeysFile string `yaml:"keys_file"`
	Location string `yaml:"location"`
}

func (c Config) Tag() string {
	return "api_keys"
}

type Key struct {
	Key            string   `yaml:"key"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

func (t *Key) Allowed(method string) bool {
	if len(t.AllowedMethods) == 0 {
		return true
	}
	for _, m := range t.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}
