# OAuth Plugin

Provides OAuth 2.0 authentication for endpoints.

## Type
- Wrapper
- Swaggerer

## Description
Implements OAuth 2.0 authentication with support for various providers (Google, GitHub, etc.). Validates access tokens and manages method-level access permissions.

## Authentication Flow

1. **Initial Setup**
   - Configure the OAuth plugin with provider credentials
   - Register your application with the OAuth provider (Google, GitHub, etc.)
   - Set up redirect URL in both plugin config and provider settings

2. **Authorization Flow**
   - Client initiates auth by accessing: `/oauth/authorize?provider=google`
   - Gateway redirects to provider's consent page
   - After user consent, provider redirects back to gateway's callback URL
   - Gateway exchanges the code for access token
   - Returns token to client for future requests

3. **Protected Endpoint Access**
   - Client includes token in requests: `Authorization: Bearer <token>`
   - Gateway validates token with provider
   - If valid and has required scopes, request proceeds
   - If invalid or insufficient scopes, returns 401/403

## Configuration

```yaml
oauth:
  provider: "google"           # OAuth provider (google, github, etc.)
  client_id: "xxx"            # OAuth Client ID
  client_secret: "xxx"        # OAuth Client Secret
  redirect_url: "http://localhost:8080/oauth/callback"  # OAuth callback URL
  auth_url: "/oauth/authorize" # Gateway's authorization endpoint (optional)
  callback_url: "/oauth/callback" # Gateway's callback endpoint (optional)
  scopes:                     # Required access scopes
    - "profile"
    - "email"
  method_scopes:             # Required scopes for specific methods (optional)
    get_users: ["read:user"]
    create_user: ["write:user"]
  token_header: "Authorization" # Header name for the token
```

## Built-in Endpoints

The plugin automatically adds these endpoints to your gateway:

### 1. Authorization Endpoint
```
GET /oauth/authorize?provider=<provider_name>
```
Initiates the OAuth flow by redirecting to the provider's consent page.

### 2. Callback Endpoint
```
GET /oauth/callback
```
Handles the OAuth provider's redirect callback:
- Exchanges authorization code for access token
- Returns token to client
- Optionally redirects to specified URL with token

### Example Usage

1. Start authentication:
```bash
# Redirect user to authorization URL
GET http://your-gateway/oauth/authorize?provider=google
```

2. After successful authentication, use the token:
```bash
# Access protected endpoint
curl -H "Authorization: Bearer <token>" http://your-gateway/api/protected-endpoint
```

## Provider-Specific Configuration

### Google
```yaml
oauth:
  provider: "google"
  client_id: "xxx.apps.googleusercontent.com"
  client_secret: "xxx"
  scopes:
    - "https://www.googleapis.com/auth/userinfo.profile"
    - "https://www.googleapis.com/auth/userinfo.email"
```

### GitHub
```yaml
oauth:
  provider: "github"
  client_id: "xxx"
  client_secret: "xxx"
  scopes:
    - "read:user"
    - "user:email"
```

## Security Considerations

1. Always use HTTPS in production
2. Keep client_secret secure
3. Validate redirect_urls to prevent open redirect vulnerabilities
4. Implement state parameter to prevent CSRF attacks
5. Validate tokens on every request
6. Use appropriate scopes - principle of least privilege 