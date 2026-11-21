package sslmonitoring

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_sslmonitoring "github.com/cancom/terraform-provider-cancom/client/services/ssl-monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceSslMonitor_AutoScanOnlyOnCreateAndRefreshesState(t *testing.T) {
	scanCalls := 0
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/Services":
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "ssl-monitoring", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
		case r.Method == "POST" && r.URL.Path == "/v1/SslMonitors":
			_ = json.NewEncoder(w).Encode(client_sslmonitoring.SslMonitor{ID: "mon-1"})
		case r.Method == "POST" && r.URL.Path == "/v1/SslMonitors/mon-1/startScan":
			scanCalls++
			w.WriteHeader(http.StatusOK)
		case r.Method == "PUT" && r.URL.Path == "/v1/SslMonitors/mon-1":
			_ = json.NewEncoder(w).Encode(client_sslmonitoring.SslMonitor{ID: "mon-1", LastUpdatedAt: 200})
		case r.Method == "GET" && r.URL.Path == "/v1/SslMonitors/mon-1":
			_ = json.NewEncoder(w).Encode(client_sslmonitoring.SslMonitor{
				ID:                 "mon-1",
				DomainName:         "example.com",
				Tenant:             "t1",
				MinimumGrade:       "A",
				ContactEmailCancom: "ops@example.com",
				Protocol:           "http",
				State:              "finished",
				CreatedBy:          "creator",
				LastUpdatedBy:      "updater",
				LastUpdatedAt:      200,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	host := srv.URL
	token := "tok"
	ccp, err := client.NewClient(&host, &token, "")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	r := resourceSslMonitor()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"domain_name":          "example.com",
		"tenant":               "t1",
		"minimum_grade":        "A",
		"contact_email_cancom": "ops@example.com",
		"protocol":             "http",
		"auto_scan":            true,
	})

	if diags := resourceSslMonitorCreate(context.Background(), d, ccp); diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if scanCalls != 1 {
		t.Fatalf("expected 1 scan call on Create, got %d", scanCalls)
	}
	if d.Get("state").(string) != "finished" || d.Get("created_by").(string) != "creator" {
		t.Fatalf("expected Create to refresh computed attributes via Read, got state=%q created_by=%q", d.Get("state"), d.Get("created_by"))
	}

	if diags := resourceSslMonitorUpdate(context.Background(), d, ccp); diags.HasError() {
		t.Fatalf("Update failed: %v", diags)
	}
	if scanCalls != 1 {
		t.Fatalf("expected Update NOT to trigger StartSslScan (total scans=%d)", scanCalls)
	}
	if d.Get("last_updated_by").(string) != "updater" {
		t.Fatalf("expected Update to refresh computed attributes via Read, got last_updated_by=%q", d.Get("last_updated_by"))
	}
}
