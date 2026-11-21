package ipam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
	client_ipam "github.com/cancom/terraform-provider-cancom/client/services/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestInstanceUpgradeV0_HandlesEmptyReleaseWaitTime(t *testing.T) {
	upgrader := instanceUpgradeV0()

	raw := map[string]interface{}{
		"id":                "inst-1",
		"name_tag":          "test",
		"managed_by":        "CUSTOMER",
		"release_wait_time": "",
		"created_at":        "2025-07-01T12:00:00",
		"updated_at":        "",
	}

	got, err := upgrader.Upgrade(context.Background(), raw, nil)
	if err != nil {
		t.Fatalf("Upgrade failed for empty release_wait_time: %v", err)
	}
	if got["release_wait_time"] != 0 {
		t.Fatalf("expected release_wait_time=0, got %#v", got["release_wait_time"])
	}
	if got["updated_at"] != nil {
		t.Fatalf("expected updated_at=nil, got %#v", got["updated_at"])
	}

	rawNum := map[string]interface{}{
		"release_wait_time": "350",
	}
	gotNum, err := upgrader.Upgrade(context.Background(), rawNum, nil)
	if err != nil {
		t.Fatalf("Upgrade failed for numeric release_wait_time: %v", err)
	}
	if gotNum["release_wait_time"] != 350 {
		t.Fatalf("expected release_wait_time=350, got %#v", gotNum["release_wait_time"])
	}
}

func TestResourceNetworkCreate_HostAssignFalsePreservesNameAndDescription(t *testing.T) {
	var updateReq client_ipam.NetworkUpdateRequest
	var putCalled bool

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v1/Services":
			_ = json.NewEncoder(w).Encode([]client.Service{
				{ServiceName: "ip-management", ServiceEndpoint: client.ServiceEndpoint{Backend: srv.URL}},
			})
		case r.Method == "POST" && r.URL.Path == "/v1/Networks/assign":
			_ = json.NewEncoder(w).Encode(client_ipam.Network{
				ID:          "net-1",
				SupernetId:  "sup-1",
				NameTag:     "my-net",
				Description: "my-desc",
				HostAssign:  true,
			})
		case r.Method == "PUT" && r.URL.Path == "/v1/Networks/net-1":
			putCalled = true
			_ = json.NewDecoder(r.Body).Decode(&updateReq)
			_ = json.NewEncoder(w).Encode(client_ipam.Network{
				ID:          "net-1",
				SupernetId:  "sup-1",
				NameTag:     updateReq.NameTag,
				Description: updateReq.Description,
				HostAssign:  updateReq.HostAssign,
			})
		case r.Method == "GET" && r.URL.Path == "/v1/Networks/net-1":
			_ = json.NewEncoder(w).Encode(client_ipam.Network{
				ID:          "net-1",
				SupernetId:  "sup-1",
				NameTag:     updateReq.NameTag,
				Description: updateReq.Description,
				HostAssign:  updateReq.HostAssign,
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

	r := resourceNetwork()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"supernet_id": "sup-1",
		"request":     "/24",
		"name_tag":    "my-net",
		"description": "my-desc",
		"host_assign": false,
	})

	diags := resourceNetworkCreate(context.Background(), d, ccp)
	if diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if !putCalled {
		t.Fatal("expected UpdateNetwork PUT call when host_assign is false")
	}
	if updateReq.NameTag != "my-net" || updateReq.Description != "my-desc" || updateReq.HostAssign != false {
		t.Fatalf("expected NameTag and Description to be preserved on PUT, got %+v", updateReq)
	}
}
