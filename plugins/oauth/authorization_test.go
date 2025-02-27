package oauth

import (
	"testing"

	"github.com/centralmind/gateway/errors"
	"github.com/stretchr/testify/require"
)

func TestAuthorization(t *testing.T) {
	// Setup test connector with config
	cfg := Config{
		AuthorizationRules: []AuthorizationRule{
			{
				// Public health check endpoints
				Methods:     []string{"GetHealth", "GetMetrics"},
				AllowPublic: true,
			},
			{
				// Only site admins can access admin endpoints
				Methods:          []string{"AdminGetUsers", "AdminDeleteUser"},
				RequireAllClaims: true,
				ClaimRules: []ClaimRule{
					{
						Claim:     "site_admin",
						Operation: "eq",
						Value:     "true",
					},
				},
			},
			{
				// Company email domain check
				Methods: []string{"InternalAPI"},
				ClaimRules: []ClaimRule{
					{
						Claim:     "email",
						Operation: "regex",
						Value:     "@double\\.cloud$",
					},
				},
			},
			{
				// Organization member check
				Methods:          []string{"GetOrganizationData"},
				RequireAllClaims: true,
				ClaimRules: []ClaimRule{
					{
						Claim:     "company",
						Operation: "eq",
						Value:     "{{.OrganizationName}}",
					},
					{
						Claim:     "type",
						Operation: "eq",
						Value:     "User",
					},
				},
			},
			{
				// Location-based access
				Methods: []string{"GetLocalData"},
				ClaimRules: []ClaimRule{
					{
						Claim:     "location",
						Operation: "eq",
						Value:     "Berlin",
					},
				},
			},
		},
	}

	connector := &Connector{config: cfg}

	// GitHub user data
	githubUser := map[string]interface{}{
		"login":      "laskoviymishka",
		"id":         2754361,
		"type":       "User",
		"site_admin": false,
		"name":       "Andrei Tserakhau",
		"company":    "Double.Cloud",
		"location":   "Berlin",
		"email":      "tserakhau@double.cloud",
		"created_at": "2012-11-08T22:37:17Z",
		"updated_at": "2025-02-18T11:17:49Z",
	}

	tests := []struct {
		name    string
		method  string
		params  map[string]interface{}
		claims  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Public endpoint should be accessible without claims",
			method:  "GetHealth",
			claims:  nil,
			wantErr: false,
		},
		{
			name:    "Admin endpoint should be denied for non-admin",
			method:  "AdminGetUsers",
			claims:  githubUser,
			wantErr: true,
		},
		{
			name:    "Should allow access for matching email domain",
			method:  "InternalAPI",
			claims:  githubUser,
			wantErr: false,
		},
		{
			name:   "Should allow access for matching organization",
			method: "GetOrganizationData",
			params: map[string]interface{}{
				"OrganizationName": "Double.Cloud",
			},
			claims:  githubUser,
			wantErr: false,
		},
		{
			name:   "Should deny access for non-matching organization",
			method: "GetOrganizationData",
			params: map[string]interface{}{
				"OrganizationName": "AnotherCompany",
			},
			claims:  githubUser,
			wantErr: true,
		},
		{
			name:    "Should allow access based on location",
			method:  "GetLocalData",
			claims:  githubUser,
			wantErr: false,
		},
		{
			name:   "Should handle non-existent claims gracefully",
			method: "GetOrganizationData",
			params: map[string]interface{}{
				"OrganizationName": "Double.Cloud",
			},
			claims: map[string]interface{}{
				"login": "laskoviymishka",
			},
			wantErr: true,
		},
		{
			name:   "Should handle template parameters correctly",
			method: "GetOrganizationData",
			params: map[string]interface{}{
				"OrganizationName": "Double.Cloud",
				"UserID":           "2754361",
			},
			claims:  githubUser,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := connector.checkAuthorization(tt.method, tt.claims, tt.params)
			if tt.wantErr {
				require.Error(t, err, "Expected error but got none")
			} else {
				require.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}

func TestClaimRuleEvaluation(t *testing.T) {
	claims := map[string]interface{}{
		"nested": map[string]interface{}{
			"array": []interface{}{
				map[string]interface{}{
					"name": "test",
					"id":   123,
				},
			},
		},
		"tags":  []string{"admin", "user"},
		"roles": []interface{}{"editor", "viewer"},
	}

	tests := []struct {
		name    string
		rule    ClaimRule
		params  map[string]interface{}
		want    bool
		wantErr bool
	}{
		{
			name: "Should access nested array value",
			rule: ClaimRule{
				Claim:     "nested.array[0].name",
				Operation: "eq",
				Value:     "test",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should handle string array contains",
			rule: ClaimRule{
				Claim:     "tags",
				Operation: "contains",
				Value:     "admin",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should handle interface array contains",
			rule: ClaimRule{
				Claim:     "roles",
				Operation: "contains",
				Value:     "editor",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should handle non-existent path",
			rule: ClaimRule{
				Claim:     "nested.nonexistent[0].field",
				Operation: "eq",
				Value:     "test",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Should handle invalid array index",
			rule: ClaimRule{
				Claim:     "nested.array[999].name",
				Operation: "eq",
				Value:     "test",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Should validate regex pattern",
			rule: ClaimRule{
				Claim:     "nested.array[0].name",
				Operation: "regex",
				Value:     "^t.*t$",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should handle invalid regex pattern",
			rule: ClaimRule{
				Claim:     "nested.array[0].name",
				Operation: "regex",
				Value:     "[invalid",
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := evaluateClaimRule(tt.rule, claims, tt.params)
			if tt.wantErr {
				require.Error(t, err, "Expected error but got none")
				return
			}

			require.NoError(t, err, "Expected no error but got: %v", err)
			require.Equal(t, tt.want, got, "Expected result %v but got %v", tt.want, got)
		})
	}
}

func TestBerlinResidencyAuthorization(t *testing.T) {
	// Setup test connector with config that makes all methods public by default
	// except BerlinCommunityEvent which requires Berlin location
	cfg := Config{
		AuthorizationRules: []AuthorizationRule{
			{
				// Default rule - allow all methods
				Methods:     []string{"*"},
				AllowPublic: true,
			},
			{
				// Override for Berlin-only community event registration
				Methods:          []string{"RegisterForBerlinCommunityEvent"},
				RequireAllClaims: true,
				ClaimRules: []ClaimRule{
					{
						Claim:     "location",
						Operation: "eq",
						Value:     "Berlin",
					},
				},
			},
		},
	}

	connector := &Connector{config: cfg}

	// Test users with different locations
	berlinUser := map[string]interface{}{
		"login":      "laskoviymishka",
		"id":         2754361,
		"type":       "User",
		"name":       "Andrei Tserakhau",
		"company":    "Double.Cloud",
		"location":   "Berlin",
		"email":      "tserakhau@double.cloud",
		"created_at": "2012-11-08T22:37:17Z",
		"updated_at": "2025-02-18T11:17:49Z",
	}

	londonUser := map[string]interface{}{
		"login":    "johndoe",
		"id":       12345678,
		"type":     "User",
		"name":     "John Doe",
		"company":  "TechCo",
		"location": "London",
		"email":    "john@techco.com",
	}

	tests := []struct {
		name    string
		method  string
		claims  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Public API should be accessible for Berlin user",
			method:  "GetPublicData",
			claims:  berlinUser,
			wantErr: false,
		},
		{
			name:    "Public API should be accessible for London user",
			method:  "GetPublicData",
			claims:  londonUser,
			wantErr: false,
		},
		{
			name:    "Public API should be accessible without authentication",
			method:  "GetPublicData",
			claims:  nil,
			wantErr: false,
		},
		{
			name:    "Berlin event registration should be allowed for Berlin user",
			method:  "RegisterForBerlinCommunityEvent",
			claims:  berlinUser,
			wantErr: false,
		},
		{
			name:    "Berlin event registration should be denied for London user",
			method:  "RegisterForBerlinCommunityEvent",
			claims:  londonUser,
			wantErr: true,
		},
		{
			name:    "Berlin event registration should be denied for unauthenticated user",
			method:  "RegisterForBerlinCommunityEvent",
			claims:  nil,
			wantErr: true,
		},
		{
			name:   "Berlin event registration should be denied for user without location",
			method: "RegisterForBerlinCommunityEvent",
			claims: map[string]interface{}{
				"login": "anonymous",
				"name":  "Anonymous User",
			},
			wantErr: true,
		},
		{
			name:    "User profile update should be accessible for any authenticated user",
			method:  "UpdateUserProfile",
			claims:  londonUser,
			wantErr: false,
		},
		{
			name:    "Health check should be accessible for everyone",
			method:  "GetHealth",
			claims:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := connector.checkAuthorization(tt.method, tt.claims, nil)
			if tt.wantErr {
				require.Error(t, err, "Expected error but got none")
				require.ErrorIs(t, err, errors.ErrNotAuthorized, "Expected unauthorized error")
			} else {
				require.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}
