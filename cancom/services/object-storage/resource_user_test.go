package object_storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_object_storage "github.com/cancom/terraform-provider-cancom/client/services/object-storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPermissionValidator_NullStatementDoesNotPanic(t *testing.T) {
	diags := permissionValidator(`{"Statement": null}`, cty.Path{})
	if !diags.HasError() {
		t.Fatal("expected error for null Statement")
	}

	valid := `{"Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`
	if diags := permissionValidator(valid, cty.Path{}); diags.HasError() {
		t.Fatalf("expected valid policy to pass, got: %v", diags)
	}
}

func TestSuppressEquivalentJSONDiffs(t *testing.T) {
	a := `{"Statement":[{"Effect":"Allow","Action":"*"}]}`
	b := "{\n  \"Statement\": [\n    {\n      \"Action\": \"*\",\n      \"Effect\": \"Allow\"\n    }\n  ]\n}"
	c := `{"Statement":[{"Effect":"Deny","Action":"*"}]}`

	if !suppressEquivalentJSONDiffs("permissions", a, b, nil) {
		t.Error("expected semantically identical JSON to suppress diff")
	}
	if suppressEquivalentJSONDiffs("permissions", a, c, nil) {
		t.Error("expected semantically different JSON not to suppress diff")
	}
}

func TestResourceUserCRUD_WritesPermissionsStringToState(t *testing.T) {
	apiPerms := map[string]interface{}{
		"Statement": []interface{}{
			map[string]interface{}{"Effect": "Allow", "Action": "s3:GetObject", "Resource": "*"},
		},
	}

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/Services" {
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "object-storage", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(client_object_storage.User{
			Id:          "u-1",
			Username:    "svc-user",
			Description: "desc",
			Permissions: apiPerms,
		})
	}))
	defer srv.Close()

	host := srv.URL
	token := "tok"
	ccp, err := client.NewClient(&host, &token, "")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	r := resourceUser()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"username":    "svc-user",
		"description": "desc",
		"permissions": `{"Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`,
	})

	if diags := resourceUserCreate(context.Background(), d, ccp); diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if got := d.Get("permissions").(string); !suppressEquivalentJSONDiffs("permissions", got, `{"Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`, nil) {
		t.Fatalf("Create did not persist API permissions JSON string in state: got %q", got)
	}

	// Simulate drift and ensure Read overwrites state with API permissions string
	_ = d.Set("permissions", `{"Statement":[]}`)
	if diags := resourceUserRead(context.Background(), d, ccp); diags.HasError() {
		t.Fatalf("Read failed: %v", diags)
	}
	if got := d.Get("permissions").(string); !suppressEquivalentJSONDiffs("permissions", got, `{"Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`, nil) {
		t.Fatalf("Read did not update permissions string in state: got %q", got)
	}

	if diags := resourceUserUpdate(context.Background(), d, ccp); diags.HasError() {
		t.Fatalf("Update failed: %v", diags)
	}
	if got := d.Get("permissions").(string); !suppressEquivalentJSONDiffs("permissions", got, `{"Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`, nil) {
		t.Fatalf("Update did not persist permissions string in state: got %q", got)
	}
}
