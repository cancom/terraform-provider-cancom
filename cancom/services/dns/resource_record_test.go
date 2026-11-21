package dns

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dns "github.com/cancom/terraform-provider-cancom/client/services/dns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newMockCcpClient(t *testing.T, handler http.HandlerFunc) (*client.CcpClient, func()) {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/Services" {
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "domdns", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
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

func TestResourceRecord_CRUDAndNotFoundRecovery(t *testing.T) {
	record := client_dns.Record{
		ID:             "rec-1",
		ZoneID:         "zone-1",
		ZoneName:       "example.com",
		Name:           "app.example.com",
		Type:           "A",
		Content:        "192.0.2.1",
		TTL:            300,
		LastChangeDate: "2026-01-01",
	}
	deleted := false

	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/v1/Records":
			_ = json.NewEncoder(w).Encode(record)
		case r.Method == "GET" && r.URL.Path == "/v1/Records/rec-1":
			if deleted {
				http.NotFound(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(record)
		case r.Method == "PUT" && r.URL.Path == "/v1/Records/rec-1":
			record.Content = "192.0.2.2"
			_ = json.NewEncoder(w).Encode(record)
		case r.Method == "DELETE" && r.URL.Path == "/v1/Records/rec-1":
			if deleted {
				http.NotFound(w, r)
				return
			}
			deleted = true
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	})
	defer cleanup()

	r := resourceRecord()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":      "app",
		"type":      "A",
		"content":   "192.0.2.1",
		"ttl":       300,
		"zone_name": "example.com",
	})

	// Create
	if diags := resourceRecordCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if d.Id() != "rec-1" {
		t.Fatalf("expected id rec-1, got %q", d.Id())
	}
	if d.Get("zone_id").(string) != "zone-1" {
		t.Fatalf("expected zone_id zone-1, got %q", d.Get("zone_id"))
	}

	// Update
	_ = d.Set("content", "192.0.2.2")
	if diags := resourceRecordUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("Update failed: %v", diags)
	}
	if d.Get("content").(string) != "192.0.2.2" {
		t.Fatalf("expected updated content 192.0.2.2, got %q", d.Get("content"))
	}

	// Delete
	if diags := resourceRecordDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("Delete failed: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected empty id after delete, got %q", d.Id())
	}

	// Read on already-deleted record clears ID without error
	d.SetId("rec-1")
	if diags := resourceRecordRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("Read on 404 should not error: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected empty id after 404 Read, got %q", d.Id())
	}

	// Delete on already-deleted record succeeds idempotently
	d.SetId("rec-1")
	if diags := resourceRecordDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("Delete on 404 should not error: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected empty id after 404 Delete, got %q", d.Id())
	}
}
