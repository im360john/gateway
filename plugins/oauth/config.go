package oauth

import "golang.org/x/oauth2"

// Config represents OAuth plugin configuration
type Config struct {
	// Provider specifies the OAuth provider ("google", "github", etc.)
	Provider string `yaml:"provider"`

	// ClientID is the OAuth Client ID
	ClientID string `yaml:"client_id"`

	// ClientSecret is the OAuth Client Secret
	ClientSecret string `yaml:"client_secret"`

	// RedirectURL for OAuth flow
	RedirectURL string `yaml:"redirect_url"`

	// Scopes defines required access scopes
	Scopes []string `yaml:"scopes"`

	// MethodScopes defines required scopes for specific methods
	MethodScopes map[string][]string `yaml:"method_scopes"`

	// TokenHeader defines the header name for the token
	TokenHeader string `yaml:"token_header"`

	// AuthURL is the gateway's authorization endpoint path
	AuthURL string `yaml:"auth_url"`

	// CallbackURL is the gateway's callback endpoint path
	CallbackURL string `yaml:"callback_url"`
}

func (c Config) Tag() string {
	return "oauth"
}

func (c Config) Doc() string {
	return docString
}

// GetOAuthConfig returns oauth2.Config for the specified provider
func (c Config) GetOAuthConfig() *oauth2.Config {
	var endpoint oauth2.Endpoint

	switch c.Provider {
	case "google":
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		}
	case "github":
		endpoint = oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		}
		// Add other providers as needed
	}

	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       c.Scopes,
		Endpoint:     endpoint,
	}
}

// WithDefaults sets default values for the config fields
func (c *Config) WithDefaults() {
	if c.TokenHeader == "" {
		c.TokenHeader = "Authorization"
	}
	if c.AuthURL == "" {
		c.AuthURL = "/oauth/authorize"
	}
	if c.CallbackURL == "" {
		c.CallbackURL = "/oauth/callback"
	}
}
