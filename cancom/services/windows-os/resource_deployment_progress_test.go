package windowsos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceWindowsOSDeploymentProgress_SchemaAndCancellation(t *testing.T) {
	r := resourceWindowsOSDeploymentProgress()

	if !r.Schema["deployment_id"].ForceNew {
		t.Error("expected deployment_id to have ForceNew: true")
	}

	validate := r.Schema["errorhandling"].ValidateDiagFunc
	if diags := validate("my_error_value", cty.Path{}); !diags.HasError() {
		t.Error("expected unanchored value 'my_error_value' to fail validation")
	}
	if diags := validate("error", cty.Path{}); diags.HasError() {
		t.Errorf("expected 'error' to pass validation, got %v", diags)
	}
	if diags := validate("warning", cty.Path{}); diags.HasError() {
		t.Errorf("expected 'warning' to pass validation, got %v", diags)
	}

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/Services" {
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "managed-windows", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
			return
		}
		// Status 1 = in progress forever
		_ = json.NewEncoder(w).Encode(client_windowsos.WindowsOS_Deplyoment{Id: "dep-1", Status: 1})
	}))
	defer srv.Close()

	host := srv.URL
	token := "tok"
	ccp, err := client.NewClient(&host, &token, "")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"deployment_id": "dep-1",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	diags := resourceWindowsOSDeploymentProgressCreate(ctx, d, ccp)
	elapsed := time.Since(start)

	if !diags.HasError() || !strings.Contains(diags[0].Summary, context.DeadlineExceeded.Error()) {
		t.Fatalf("expected context deadline exceeded diagnostic, got %v", diags)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("polling loop ignored context cancellation (took %v)", elapsed)
	}
}
