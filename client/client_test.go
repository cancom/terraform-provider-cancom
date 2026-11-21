package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func makeTestJWT(exp int64) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payloadBytes, _ := json.Marshal(map[string]int64{"exp": exp})
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return header + "." + payload + ".sig"
}

func TestDefaultHostURLHasNoV1Suffix(t *testing.T) {
	if strings.HasSuffix(HostURL, "/v1") {
		t.Fatalf("HostURL %q must not end with /v1 because service_registry appends /v1/Services", HostURL)
	}
}

func TestNewClient_RoleAssumptionAndJSONEscaping(t *testing.T) {
	var assumedRole string
	var contentType string
	expectedJWT := makeTestJWT(time.Now().Add(time.Hour).Unix())

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/Services":
			_ = json.NewEncoder(w).Encode([]Service{
				{ServiceName: "iam", ServiceEndpoint: ServiceEndpoint{Backend: srv.URL}},
				{ServiceName: "dns", ServiceEndpoint: ServiceEndpoint{Backend: srv.URL}},
			})
		case "/v1/Services/iam":
			_ = json.NewEncoder(w).Encode(Service{
				ServiceName:     "iam",
				ServiceEndpoint: ServiceEndpoint{Backend: srv.URL},
			})
		case "/v1/AssumeRole":
			contentType = r.Header.Get("Content-Type")
			var body struct {
				Role string `json:"role"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			assumedRole = body.Role
			_ = json.NewEncoder(w).Encode(map[string]string{"jwt": expectedJWT})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	host := srv.URL
	token := "base-token"
	// Role containing characters that would break fmt.Sprintf JSON formatting
	role := `crn:123::iam:role:weird"role\name`

	c, err := NewClient(&host, &token, role)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	got, err := c.tp.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if got != expectedJWT {
		t.Fatalf("expected assumed JWT %q, got %q", expectedJWT, got)
	}
	if assumedRole != role {
		t.Fatalf("expected properly JSON-escaped role %q, got %q", role, assumedRole)
	}
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}
}

func TestNewClient_RoleInitFailurePropagatesError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	host := srv.URL
	token := "base-token"
	_, err := NewClient(&host, &token, "some-role")
	if err == nil {
		t.Fatal("expected error when service-registry initialize fails")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 HTTPError, got %v", err)
	}
}

func TestWithToken_UsesNewTokenAndDoesNotPanic(t *testing.T) {
	var mu sync.Mutex
	seenAuth := map[string]string{}

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/Services":
			_ = json.NewEncoder(w).Encode([]Service{
				{ServiceName: "svc-a", ServiceEndpoint: ServiceEndpoint{Backend: srv.URL}},
				{ServiceName: "svc-b", ServiceEndpoint: ServiceEndpoint{Backend: srv.URL}},
			})
		default:
			mu.Lock()
			seenAuth[r.URL.Path] = r.Header.Get("Authorization")
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	defer srv.Close()

	host := srv.URL
	origToken := "orig-token"
	c, err := NewClient(&host, &origToken, "")
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Prime svc-a in cache before cloning
	if _, err := c.GetService("svc-a"); err != nil {
		t.Fatalf("GetService(svc-a) failed: %v", err)
	}

	cloned := c.WithToken("new-token")

	// Existing cached service must use new-token
	svcA, err := cloned.GetService("svc-a")
	if err != nil {
		t.Fatalf("cloned.GetService(svc-a) failed: %v", err)
	}
	reqA, _ := http.NewRequest("GET", svcA.HostURL+"/a", nil)
	if _, err := svcA.DoRequest(reqA); err != nil {
		t.Fatalf("svcA.DoRequest failed: %v", err)
	}

	// Newly resolved service must not panic on nil tp and must use new-token
	svcB, err := cloned.GetService("svc-b")
	if err != nil {
		t.Fatalf("cloned.GetService(svc-b) failed: %v", err)
	}
	reqB, _ := http.NewRequest("GET", svcB.HostURL+"/b", nil)
	if _, err := svcB.DoRequest(reqB); err != nil {
		t.Fatalf("svcB.DoRequest failed: %v", err)
	}

	// Original client must still use orig-token
	origSvcA, _ := c.GetService("svc-a")
	reqOrig, _ := http.NewRequest("GET", origSvcA.HostURL+"/orig", nil)
	if _, err := origSvcA.DoRequest(reqOrig); err != nil {
		t.Fatalf("origSvcA.DoRequest failed: %v", err)
	}

	if seenAuth["/a"] != "Bearer new-token" {
		t.Errorf("expected cached service to send new-token, got %q", seenAuth["/a"])
	}
	if seenAuth["/b"] != "Bearer new-token" {
		t.Errorf("expected new service to send new-token, got %q", seenAuth["/b"])
	}
	if seenAuth["/orig"] != "Bearer orig-token" {
		t.Errorf("expected original client to keep orig-token, got %q", seenAuth["/orig"])
	}
}

func TestDoRequestWithRetry_RespectsContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := newHttpClient(srv.URL, &tokenProvider{mu: &sync.Mutex{}, token: "tok"})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", srv.URL+"/retry", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext failed: %v", err)
	}

	start := time.Now()
	_, err = c.DoRequestWithRetry(req, nil)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("DoRequestWithRetry ignored context cancellation (took %v)", elapsed)
	}
}
