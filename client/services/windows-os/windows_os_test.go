package client_windowsos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
)

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	token := "test-token"
	ccp, err := client.NewClient(&serverURL, &token, "")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	svc, err := ccp.GetService("managed-windows")
	if err != nil {
		t.Fatalf("GetService failed: %v", err)
	}
	return (*Client)(svc)
}

func TestUpdateWindowsOsDeployment_UsesIDAndReturnsErrors(t *testing.T) {
	var requestedGetPath string
	var putCalled bool
	var getStatusCode = http.StatusOK
	var existing = WindowsOS_Deplyoment{
		Id:                     "dep-123",
		CoustomerEnvironmentID: "ENV-1",
		Computer: WindowsOS_Computer{
			Role:                "AppServer",
			Services:            []string{"OS"},
			MaintenanceWindowId: []string{"MW-1"},
		},
	}

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/Services":
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "managed-windows", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
		case r.Method == "GET" && r.URL.Path == "/api/v1/Deployment/dep-123":
			requestedGetPath = r.URL.Path
			if getStatusCode != http.StatusOK {
				http.Error(w, "boom", getStatusCode)
				return
			}
			_ = json.NewEncoder(w).Encode(existing)
		case r.Method == "PUT" && r.URL.Path == "/api/v1/Deployment/dep-123":
			putCalled = true
			_ = json.NewEncoder(w).Encode(existing)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)

	// 1. Valid update with empty Id on payload struct (like resource_deployment.go passes)
	resp, err := c.UpdateWindowsOsDeployment("dep-123", &WindowsOS_Deplyoment{
		CoustomerEnvironmentID: "ENV-1",
		Computer: WindowsOS_Computer{
			Role:                "AppServer",
			Services:            []string{"OS"},
			MaintenanceWindowId: []string{"MW-2"},
		},
	})
	if err != nil {
		t.Fatalf("expected valid update to succeed, got err: %v", err)
	}
	if resp == nil || resp.Id != "dep-123" {
		t.Fatalf("expected non-nil response with id dep-123, got %#v", resp)
	}
	if requestedGetPath != "/api/v1/Deployment/dep-123" {
		t.Fatalf("expected GET /api/v1/Deployment/dep-123, got %q", requestedGetPath)
	}
	if !putCalled {
		t.Fatal("expected PUT to be called")
	}

	// 2. GET failure must return a non-nil error (not nil, nil)
	getStatusCode = http.StatusInternalServerError
	resp, err = c.UpdateWindowsOsDeployment("dep-123", &WindowsOS_Deplyoment{
		CoustomerEnvironmentID: "ENV-1",
		Computer:               WindowsOS_Computer{Role: "AppServer"},
	})
	if err == nil || resp != nil {
		t.Fatalf("expected (nil, non-nil err) when GET fails, got (%#v, %v)", resp, err)
	}

	// 3. Immutable field mismatch must return a non-nil error (not nil, nil)
	getStatusCode = http.StatusOK
	resp, err = c.UpdateWindowsOsDeployment("dep-123", &WindowsOS_Deplyoment{
		CoustomerEnvironmentID: "ENV-OTHER",
		Computer:               WindowsOS_Computer{Role: "AppServer"},
	})
	if err == nil || resp != nil {
		t.Fatalf("expected error on customerEnvironmentID mismatch, got (%#v, %v)", resp, err)
	}

	resp, err = c.UpdateWindowsOsDeployment("dep-123", &WindowsOS_Deplyoment{
		CoustomerEnvironmentID: "ENV-1",
		Computer:               WindowsOS_Computer{Role: "OtherRole"},
	})
	if err == nil || resp != nil {
		t.Fatalf("expected error on role mismatch, got (%#v, %v)", resp, err)
	}
}
