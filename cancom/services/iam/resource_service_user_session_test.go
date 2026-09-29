package iam

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
)

const testJWT = "eyJhbGciOiJFUzI1NiIsImtpZCI6IjI5ODU4NDkxLTAwMGQtNDFiMy05YzcxLTk2MDlhMTMxM2NmZSIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIxIn0.eyJpYXQiOjE3OTAxODQzMTUsImV4cCI6MTgyMTcyMDMxNSwiaXNzIjoiaHR0cHM6Ly9pYW0ucG9ydGFsLmNhbmNvbS5kZXYiLCJhdWQiOlsiSUFNLVByaW5jaXBhbCJdLCJzZXNzaW9uSWQiOiJhRDFfbzRDdUVZSVJ0bnVvMEdEQUo1eHJuX05QSFRDbUVyOGdNVzNibi1rTTc2cmhZRFN4dy05eW9yYjhUMzZXIiwidHlwIjoiQmVhcmVyIiwidGVuYW50IjoiY2FuY29tIiwic3ViYWNjb3VudCI6IiIsInByaW5jaXBhbFR5cGUiOiJzZXJ2aWNldXNlciIsInByaW5jaXBhbE5hbWUiOiJ0ZXN0dXNlcjIiLCJwcmluY2lwYWxDUk4iOiJjcm46Y2FuY29tOjppYW06c2VydmljZXVzZXI6dGVzdHVzZXIyIn0.cr7gHMqm_jykX2LjooLzuk_d3DHfuDHlgGTPWXuhfvTm6F8-Ct2DtnXwK0xw90Qn5BW8Y1wmCl2Iamzy-uBhsw"

func TestResourceServiceUserSessionCRUD(t *testing.T) {
	principalCRN := "crn:cancom::iam:serviceuser:testuser2"
	sessionID := "aD1_o4CuEYIRtnuo0GDAJ5xrn_NPHTCmEr8gMW3bn-kM76rhYDSxw-9yorb8T36W"
	var currentTTL int64 = time.Now().Add(30 * 24 * time.Hour).Unix()
	var currentComment = "test initial comment"

	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/Services":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"serviceName": "iam",
					"serviceEndpoint": map[string]string{
						"backend": serverURL,
					},
				},
			})

		case r.Method == http.MethodPost && r.URL.Path == "/v1/Sessions":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"jwt": testJWT,
			})

		case r.Method == http.MethodGet && r.URL.Path == "/v1/Sessions/"+principalCRN:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"enabled":      true,
					"lastIssuedAt": strconv.FormatInt(time.Now().Unix(), 10),
					"comment":      currentComment,
					"ttl":          strconv.FormatInt(currentTTL, 10),
					"principal":    principalCRN,
					"sessionId":    sessionID,
				},
			})

		case r.Method == http.MethodPut && r.URL.Path == "/v1/Sessions/"+principalCRN+"/"+sessionID:
			var req map[string]string
			_ = json.NewDecoder(r.Body).Decode(&req)
			currentComment = req["comment"]
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodDelete && r.URL.Path == "/v1/Sessions/"+principalCRN+"/"+sessionID:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "deleted",
			})

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	token := "dummy-token"
	ccpClient, err := client.NewClient(&server.URL, &token, "")
	if err != nil {
		t.Fatalf("failed to create ccpClient: %v", err)
	}

	r := resourceServiceUserSession()
	d := r.TestResourceData()
	_ = d.Set("service_user", principalCRN)
	_ = d.Set("comment", "test initial comment")
	_ = d.Set("reroll_days", 7)

	ctx := context.Background()

	// 1. Test Create
	diags := resourceServiceUserSessionCreate(ctx, d, ccpClient)
	if diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if d.Id() != sessionID {
		t.Errorf("Expected ID %s, got %s", sessionID, d.Id())
	}
	if d.Get("jwt") != testJWT {
		t.Errorf("Expected JWT %s, got %s", testJWT, d.Get("jwt"))
	}
	if d.Get("session_id") != sessionID {
		t.Errorf("Expected session_id %s, got %s", sessionID, d.Get("session_id"))
	}

	// 2. Test Read when valid (TTL is 30 days away, reroll_days is 7)
	diags = resourceServiceUserSessionRead(ctx, d, ccpClient)
	if diags.HasError() {
		t.Fatalf("Read failed: %v", diags)
	}
	if d.Id() != sessionID {
		t.Errorf("Resource ID should remain set, got: %s", d.Id())
	}

	// 3. Test Update comment
	_ = d.Set("comment", "updated comment")
	diags = resourceServiceUserSessionUpdate(ctx, d, ccpClient)
	if diags.HasError() {
		t.Fatalf("Update failed: %v", diags)
	}
	if currentComment != "updated comment" {
		t.Errorf("Expected updated comment on server, got %s", currentComment)
	}

	// 4. Read refreshes expires_at from the API TTL
	currentTTL = time.Now().Add(3 * 24 * time.Hour).Unix()
	diags = resourceServiceUserSessionRead(ctx, d, ccpClient)
	if diags.HasError() {
		t.Fatalf("Read failed: %v", diags)
	}
	if int64(d.Get("expires_at").(int)) != currentTTL {
		t.Errorf("Expected expires_at %d, got %d", currentTTL, d.Get("expires_at").(int))
	}

	// 5. Test Delete
	d.SetId(sessionID)
	diags = resourceServiceUserSessionDelete(ctx, d, ccpClient)
	if diags.HasError() {
		t.Fatalf("Delete failed: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("Expected ID to be cleared after delete, got %s", d.Id())
	}
}

func TestSessionNeedsReroll(t *testing.T) {
	now := time.Now()
	day := int64(24 * 60 * 60)
	cases := []struct {
		name       string
		expiresAt  int64
		rerollDays int
		want       bool
	}{
		{"unknown expiry", 0, 7, false},
		{"far from expiry", now.Unix() + 30*day, 7, false},
		{"within reroll window", now.Unix() + 3*day, 7, true},
		{"already expired", now.Unix() - day, 0, true},
		{"reroll 0 and still valid", now.Unix() + day, 0, false},
	}
	for _, tc := range cases {
		if got := sessionNeedsReroll(tc.expiresAt, tc.rerollDays, now); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestProviderResourceRegistered(t *testing.T) {
	p := New()
	if _, ok := p.Provider().ResourcesMap["service_user_session"]; !ok {
		t.Errorf("expected service_user_session to be registered in provider ResourcesMap")
	}
}
