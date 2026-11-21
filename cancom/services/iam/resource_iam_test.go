package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newMockCcpClient(t *testing.T, handler http.HandlerFunc) (*client.CcpClient, func()) {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/Services" {
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "iam", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
			return
		}
		handler(w, r)
	}))

	host := srv.URL
	token := "test-token"
	c, err := client.NewClient(&host, &token, "")
	if err != nil {
		srv.Close()
		t.Fatalf("NewClient failed: %v", err)
	}
	return c, srv.Close
}

func TestResourceUserCreate_SetsIDBeforeRead(t *testing.T) {
	var getPath string
	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/v1/Users":
			_ = json.NewEncoder(w).Encode(client_iam.UserCreateResponse{Principal: "crn:t::iam:user:alice"})
		case r.Method == "GET":
			getPath = r.URL.Path
			_ = json.NewEncoder(w).Encode(client_iam.User{
				Principal: "crn:t::iam:user:alice",
				Tenant:    "t",
				CreatedAt: "2026-01-01",
				CreatedBy: "admin",
			})
		}
	})
	defer cleanup()

	r := resourceUser()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "alice"})
	diags := resourceUserCreate(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if getPath != "/v1/Users/crn:t::iam:user:alice" {
		t.Fatalf("expected Read to fetch created principal, got %q", getPath)
	}
	if d.Get("tenant").(string) != "t" {
		t.Fatalf("expected tenant to be populated by Read, got %q", d.Get("tenant"))
	}
}

func TestResourceRoleCreate_SetsIDBeforeRead(t *testing.T) {
	var getPath string
	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/v1/Roles":
			_ = json.NewEncoder(w).Encode(client_iam.RoleCreateResponse{Principal: "crn:t::iam:role:dev"})
		case r.Method == "GET":
			getPath = r.URL.Path
			_ = json.NewEncoder(w).Encode(client_iam.Role{
				Principal: "crn:t::iam:role:dev",
				Tenant:    "t",
				CreatedAt: "2026-01-01",
				CreatedBy: "admin",
			})
		}
	})
	defer cleanup()

	r := resourceRole()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "dev"})
	diags := resourceRoleCreate(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if getPath != "/v1/Roles/crn:t::iam:role:dev" {
		t.Fatalf("expected Read to fetch created principal, got %q", getPath)
	}
	if d.Get("tenant").(string) != "t" {
		t.Fatalf("expected tenant to be populated by Read, got %q", d.Get("tenant"))
	}
}

func TestResourceServiceUserCreate_SensitiveSessionAndNoWarningLeak(t *testing.T) {
	secretSession := "super-secret-session-token"
	var getPath string
	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/v1/ServiceUsers":
			_ = json.NewEncoder(w).Encode(client_iam.ServiceUserCreateResponse{
				Principal: "crn:t::iam:service-user:bot",
				Session:   secretSession,
			})
		case r.Method == "GET":
			getPath = r.URL.Path
			_ = json.NewEncoder(w).Encode(client_iam.ServiceUser{
				Principal: "crn:t::iam:service-user:bot",
				Tenant:    "t",
				CreatedAt: "2026-01-01",
				CreatedBy: "admin",
			})
		}
	})
	defer cleanup()

	r := resourceServiceUser()
	if !r.Schema["session"].Sensitive {
		t.Fatal("expected session attribute to be marked Sensitive")
	}

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"name": "bot"})
	diags := resourceServiceUserCreate(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	for _, diagItem := range diags {
		if strings.Contains(diagItem.Detail, secretSession) || strings.Contains(diagItem.Summary, secretSession) {
			t.Fatalf("secret session token leaked in diagnostics: %+v", diagItem)
		}
	}
	if getPath != "/v1/ServiceUsers/crn:t::iam:service-user:bot" {
		t.Fatalf("expected Read to fetch created principal, got %q", getPath)
	}
	if d.Get("session").(string) != secretSession {
		t.Fatalf("expected session state to be set, got %q", d.Get("session"))
	}
	if d.Get("tenant").(string) != "t" {
		t.Fatalf("expected tenant to be populated by Read, got %q", d.Get("tenant"))
	}
}

func TestResourcePolicyCreate_PrecedenceAndValidation(t *testing.T) {
	var sentPolicy client_iam.PolicyRequest
	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&sentPolicy)
		w.WriteHeader(http.StatusOK)
	})
	defer cleanup()

	r := resourcePolicy()

	// 1. Custom takes precedence over profile when both are specified
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"service":   "domdns",
		"principal": "crn:t::iam:user:alice",
		"policy": []interface{}{
			map[string]interface{}{
				"profile": "full-access",
				"custom": map[string]interface{}{
					"listRecords": "*",
				},
			},
		},
	})
	diags := resourcePolicyCreate(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if len(sentPolicy.Policy) != 1 || sentPolicy.Policy["listRecords"] != "*" {
		t.Fatalf("expected custom policy to override profile, got %#v", sentPolicy.Policy)
	}

	// 2. Empty policy list returns error without panic
	dEmptyList := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"service":   "domdns",
		"principal": "crn:t::iam:user:alice",
		"policy":    []interface{}{},
	})
	if diags := resourcePolicyCreate(context.Background(), dEmptyList, c); !diags.HasError() {
		t.Fatal("expected error for empty policy list")
	}

	// 3. Neither custom nor profile returns error
	dEmptyBlock := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"service":   "domdns",
		"principal": "crn:t::iam:user:alice",
		"policy":    []interface{}{map[string]interface{}{}},
	})
	if diags := resourcePolicyCreate(context.Background(), dEmptyBlock, c); !diags.HasError() {
		t.Fatal("expected error when neither custom nor profile is set")
	}

	// 4. Unknown service with profile returns error instead of empty policy
	dUnknownSvc := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"service":   "nonexistent-service",
		"principal": "crn:t::iam:user:alice",
		"policy": []interface{}{
			map[string]interface{}{"profile": "reader"},
		},
	})
	if diags := resourcePolicyCreate(context.Background(), dUnknownSvc, c); !diags.HasError() {
		t.Fatal("expected error for unknown service profile")
	}

	// 5. Invalid profile returns error
	dBadProfile := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"service":   "domdns",
		"principal": "crn:t::iam:user:alice",
		"policy": []interface{}{
			map[string]interface{}{"profile": "full_access"},
		},
	})
	if diags := resourcePolicyCreate(context.Background(), dBadProfile, c); !diags.HasError() {
		t.Fatal("expected error for invalid profile name")
	}

	// 6. LowerInitial handles empty string without panic
	if got := LowerInitial(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
