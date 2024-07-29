package client

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const HostURL string = "https://service-registry.portal.cancom.dev/v1"

type Client struct {
	HostURL     string
	ServiceURLs map[string]string
	HTTPClient  *http.Client
	Token       string
}

func NewClient(host, token *string) (*Client, error) {
	c := Client{
		HTTPClient:  &http.Client{Timeout: 2 * time.Minute},
		ServiceURLs: map[string]string{},
		HostURL:     HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if token == nil {
		return nil, fmt.Errorf("Token is required")
	}

	c.Token = *token

	return &c, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != 201) {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
	}

	return body, nil
}

func (c *Client) DoRequestWithRetry(req *http.Request, policy func(resp *http.Response) bool) ([]byte, error) {
	token := c.Token
	if policy == nil {
		policy = func(resp *http.Response) bool {
			return resp.StatusCode == 429
		}
	}

	req.Header.Add("Authorization", "Bearer "+token)

	var bodyGetter = func() io.ReadCloser {
		return http.NoBody
	}

	if req.Body != nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		defer req.Body.Close()

		bodyGetter = func() io.ReadCloser {
			return io.NopCloser(bytes.NewReader(buf))
		}
	}

	for i := 0; i < 30; i++ {
		req.Body = bodyGetter()
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if policy(resp) {
			time.Sleep(time.Duration(rand.Intn(16)+15) * time.Second)
			continue
		}

		if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != 201) {
			return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
		}

		return body, nil
	}

	return nil, fmt.Errorf("maximum retries has been reached")
}
