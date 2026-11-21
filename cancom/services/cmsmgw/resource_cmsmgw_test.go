package cmsmgw

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_cmsmgw "github.com/cancom/terraform-provider-cancom/client/services/cmsmgw"
)

func newMockCcpClient(t *testing.T, handler http.HandlerFunc) (*client.CcpClient, func()) {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/Services" {
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "cmsmgw", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
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

func TestResourceConnection_SensitivePSKs(t *testing.T) {
	r := resourceConnection()
	for _, attr := range []string{"ipsec_psk_a", "ipsec_psk_b"} {
		if !r.Schema[attr].Sensitive {
			t.Errorf("expected %s to be marked Sensitive", attr)
		}
	}
}

func TestResourceConnectionDelete_NoPanicOn404Or500(t *testing.T) {
	for _, tc := range []struct {
		name       string
		getStatus  int
		wantErrors bool
	}{
		{name: "immediate 404 completes cleanly", getStatus: http.StatusNotFound, wantErrors: false},
		{name: "500 returns error without panic", getStatus: http.StatusInternalServerError, wantErrors: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "DELETE":
					_ = json.NewEncoder(w).Encode(client_cmsmgw.Connection{ID: "conn-1"})
				case "GET":
					http.Error(w, "status", tc.getStatus)
				}
			})
			defer cleanup()

			r := resourceConnection()
			d := r.TestResourceData()
			d.SetId("conn-1")

			diags := resourceConnectionDelete(context.Background(), d, c)
			if diags.HasError() != tc.wantErrors {
				t.Fatalf("wantErrors=%v, got diags=%v", tc.wantErrors, diags)
			}
			if !tc.wantErrors && d.Id() != "" {
				t.Fatalf("expected ID to be cleared after delete, got %q", d.Id())
			}
		})
	}
}

func TestResourceTranslationDelete_NoPanicOn404Or500(t *testing.T) {
	for _, tc := range []struct {
		name       string
		getStatus  int
		wantErrors bool
	}{
		{name: "immediate 404 completes cleanly", getStatus: http.StatusNotFound, wantErrors: false},
		{name: "500 returns error without panic", getStatus: http.StatusInternalServerError, wantErrors: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case "DELETE":
					_ = json.NewEncoder(w).Encode(client_cmsmgw.Translation{ID: "tr-1", DeploymentState: "deleted"})
				case "GET":
					http.Error(w, "status", tc.getStatus)
				}
			})
			defer cleanup()

			r := resourceTranslation()
			d := r.TestResourceData()
			d.SetId("tr-1")

			diags := resourceTranslationDelete(context.Background(), d, c)
			if diags.HasError() != tc.wantErrors {
				t.Fatalf("wantErrors=%v, got diags=%v", tc.wantErrors, diags)
			}
			if !tc.wantErrors && d.Id() != "" {
				t.Fatalf("expected ID to be cleared after delete, got %q", d.Id())
			}
		})
	}
}

func TestResourceGatewayRead_SetsBastionLiteWindows(t *testing.T) {
	c, cleanup := newMockCcpClient(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(client_cmsmgw.Gateway{
			ID:                 "gw-1",
			BastionLiteLinux:   true,
			BastionLiteWindows: true,
		})
	})
	defer cleanup()

	r := resourceGateway()
	d := r.TestResourceData()
	d.SetId("gw-1")

	diags := resourceGatewayRead(context.Background(), d, c)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if !d.Get("bastion_lite_windows").(bool) {
		t.Fatal("expected bastion_lite_windows to be true after Read")
	}
}
