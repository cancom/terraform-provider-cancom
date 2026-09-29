package client_iam

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cancom/terraform-provider-cancom/client"
)

const sampleJWT = "eyJhbGciOiJFUzI1NiIsImtpZCI6IjI5ODU4NDkxLTAwMGQtNDFiMy05YzcxLTk2MDlhMTMxM2NmZSIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIxIn0.eyJpYXQiOjE3OTAxODQzMTUsImV4cCI6MTgyMTcyMDMxNSwiaXNzIjoiaHR0cHM6Ly9pYW0ucG9ydGFsLmNhbmNvbS5kZXYiLCJhdWQiOlsiSUFNLVByaW5jaXBhbCJdLCJzZXNzaW9uSWQiOiJhRDFfbzRDdUVZSVJ0bnVvMEdEQUo1eHJuX05QSFRDbUVyOGdNVzNibi1rTTc2cmhZRFN4dy05eW9yYjhUMzZXIiwidHlwIjoiQmVhcmVyIiwidGVuYW50IjoiY2FuY29tIiwic3ViYWNjb3VudCI6IiIsInByaW5jaXBhbFR5cGUiOiJzZXJ2aWNldXNlciIsInByaW5jaXBhbE5hbWUiOiJ0ZXN0dXNlcjIiLCJwcmluY2lwYWxDUk4iOiJjcm46Y2FuY29tOjppYW06c2VydmljZXVzZXI6dGVzdHVzZXIyIn0.cr7gHMqm_jykX2LjooLzuk_d3DHfuDHlgGTPWXuhfvTm6F8-Ct2DtnXwK0xw90Qn5BW8Y1wmCl2Iamzy-uBhsw"

func TestParseSessionToken(t *testing.T) {
	claims, err := ParseSessionToken(sampleJWT)
	if err != nil {
		t.Fatalf("unexpected error parsing jwt: %v", err)
	}

	if claims.SessionID != "aD1_o4CuEYIRtnuo0GDAJ5xrn_NPHTCmEr8gMW3bn-kM76rhYDSxw-9yorb8T36W" {
		t.Errorf("unexpected SessionID: got %s", claims.SessionID)
	}

	if claims.PrincipalCRN != "crn:cancom::iam:serviceuser:testuser2" {
		t.Errorf("unexpected PrincipalCRN: got %s", claims.PrincipalCRN)
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Unix() != 1821720315 {
		t.Errorf("unexpected ExpiresAt: got %v", claims.ExpiresAt)
	}

	if claims.Tenant != "cancom" {
		t.Errorf("unexpected Tenant: got %s", claims.Tenant)
	}

	if claims.PrincipalType != "serviceuser" {
		t.Errorf("unexpected PrincipalType: got %s", claims.PrincipalType)
	}

	if claims.PrincipalName != "testuser2" {
		t.Errorf("unexpected PrincipalName: got %s", claims.PrincipalName)
	}
}

func TestSessionClientMethods(t *testing.T) {
	principalCRN := "crn:cancom::iam:serviceuser:testuser2"
	sessionID := "aD1_o4CuEYIRtnuo0GDAJ5xrn_NPHTCmEr8gMW3bn-kM76rhYDSxw-9yorb8T36W"

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
			var req SessionCreateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.ServiceUser != principalCRN || req.Comment != "test mbr" {
				http.Error(w, "unexpected request body", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"jwt": sampleJWT,
			})

		case r.Method == http.MethodGet && r.URL.Path == "/v1/Sessions/"+principalCRN:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]Session{
				{
					Enabled:      true,
					LastIssuedAt: "1790184826",
					Comment:      "test",
					TTL:          "1821720826",
					Principal:    principalCRN,
					SessionID:    sessionID,
				},
			})

		case r.Method == http.MethodPut && r.URL.Path == "/v1/Sessions/"+principalCRN+"/"+sessionID:
			var req SessionUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.Comment != "test 123" {
				http.Error(w, "unexpected request body", http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodDelete && r.URL.Path == "/v1/Sessions/"+principalCRN+"/"+sessionID:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "successfully removed session " + sessionID + " from principal " + principalCRN,
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

	rawClient, err := ccpClient.GetService("iam")
	if err != nil {
		t.Fatalf("failed to get iam service: %v", err)
	}
	iamClient := (*Client)(rawClient)

	// Test CreateSession
	createResp, err := iamClient.CreateSession(&SessionCreateRequest{
		ServiceUser: principalCRN,
		Comment:     "test mbr",
	})
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if createResp.Jwt != sampleJWT {
		t.Errorf("CreateSession: unexpected Jwt: %s", createResp.Jwt)
	}
	if createResp.SessionID != sessionID {
		t.Errorf("CreateSession: unexpected SessionID: %s", createResp.SessionID)
	}
	if createResp.PrincipalCRN != principalCRN {
		t.Errorf("CreateSession: unexpected PrincipalCRN: %s", createResp.PrincipalCRN)
	}
	if createResp.ExpiresAt != 1821720315 {
		t.Errorf("CreateSession: unexpected ExpiresAt: %d", createResp.ExpiresAt)
	}

	// Test GetSession
	session, err := iamClient.GetSession(principalCRN, sessionID)
	if err != nil {
		t.Fatalf("GetSession failed: %v", err)
	}
	if session.SessionID != sessionID || session.Principal != principalCRN {
		t.Errorf("GetSession: unexpected session returned: %+v", session)
	}

	// Test GetSession Not Found
	nonExistentSession, err := iamClient.GetSession(principalCRN, "non-existent-id")
	if err == nil {
		t.Errorf("expected error when getting non-existent session, got %v", nonExistentSession)
	}

	// Test ListSessions
	sessions, err := iamClient.ListSessions(principalCRN)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}
	if len(sessions) != 1 || sessions[0].SessionID != sessionID {
		t.Errorf("ListSessions: unexpected sessions: %+v", sessions)
	}

	// Test UpdateSession
	err = iamClient.UpdateSession(principalCRN, sessionID, &SessionUpdateRequest{
		Comment: "test 123",
	})
	if err != nil {
		t.Fatalf("UpdateSession failed: %v", err)
	}

	// Test DeleteSession
	err = iamClient.DeleteSession(principalCRN, sessionID)
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}
}
