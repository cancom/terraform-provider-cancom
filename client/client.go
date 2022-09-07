package client

import (
	"fmt"
	"io/ioutil"
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
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
	}

	return body, nil
}
