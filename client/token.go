package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tokenProvider struct {
	mu           *sync.Mutex
	token        string
	currentToken string
	expireTime   int64
	role         string

	iamClient *Client
}

func (t *tokenProvider) setToken(token string) error {

	// Parse the JWT token to check expiration
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format: expected 3 parts but got %d", len(parts))
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	// Parse the JSON payload
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return err
	}

	t.currentToken = token
	t.expireTime = claims.Exp - int64((5 * time.Minute).Seconds())

	return nil
}

func (t *tokenProvider) GetToken() (string, error) {
	// base case - no role never needs to check expired tokens
	if t.role == "" {
		return t.token, nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if time.Now().Unix() > t.expireTime {
		// need to refresh token
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/AssumeRole", t.iamClient.HostURL), strings.NewReader(fmt.Sprintf(`{"role": "%s"}`, t.role)))
		if err != nil {
			return "", err
		}

		resp, err := t.iamClient.DoRequest(req)
		if err != nil {
			return "", err
		}

		var token struct {
			Token string `json:"jwt"`
		}

		err = json.Unmarshal(resp, &token)
		if err != nil {
			return "", err
		}

		if err := t.setToken(token.Token); err != nil {
			return "", err
		}
	}

	return t.currentToken, nil
}
