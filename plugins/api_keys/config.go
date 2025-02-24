package api_keys

type Config struct {
	Name      string  `yaml:"name"`
	Tokens    []Token `yaml:"tokens"`
	TokenFile string  `yaml:"token_file"`
	Location  string  `yaml:"location"`
}

func (c Config) Tag() string {
	return "api_keys"
}

type Token struct {
	Token          string   `yaml:"token"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

func (t *Token) Allowed(method string) bool {
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
