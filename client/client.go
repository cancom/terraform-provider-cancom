package client

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const HostURL string = "https://service-registry.portal.cancom.dev/v1"

type CcpClient struct {
	client      *http.Client
	services    map[string]*Client
	serviceURLs map[string]string
	token       string
	initialized bool
	mu          sync.Mutex
}

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	token      string
}

func NewClient(host, token *string) (*CcpClient, error) {
	serviceRegistry := HostURL
	if host != nil {
		serviceRegistry = *host
	}

	if token == nil {
		return nil, fmt.Errorf("token is required")
	}

	if serviceRegistry == "" {
		return nil, fmt.Errorf("service registry is not valid, set to: %s", serviceRegistry)
	}

	c := CcpClient{
		client:      &http.Client{Timeout: 2 * time.Minute},
		services:    map[string]*Client{"service-registry": newHttpClient(serviceRegistry, *token)},
		serviceURLs: map[string]string{},
		initialized: false,
		token:       *token,
	}

	return &c, nil
}

// initialize the ccp client
func (c *CcpClient) initialize() error {
	services, err := c.getAllServices()

	if err != nil {
		return err
	}

	for _, s := range services {
		if s.ServiceEndpoint.Backend != "" {
			c.serviceURLs[s.ServiceName] = s.ServiceEndpoint.Backend
		}
	}

	c.initialized = true
	return nil
}

func newHttpClient(host, token string) *Client {
	return &Client{
		HostURL:    host,
		token:      token,
		HTTPClient: &http.Client{Timeout: 2 * time.Minute},
	}
}

func (c *CcpClient) GetService(name string) (*Client, error) {
	// Get access to lock
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		err := c.initialize()

		if err != nil {
			return nil, err
		}
	}

	serviceClient, ok := c.services[name]
	if !ok {
		serviceUrl, ok := c.serviceURLs[name]

		if !ok {
			return nil, fmt.Errorf("service %s is not known", name)
		}

		serviceClient = newHttpClient(serviceUrl, c.token)
		c.services[name] = serviceClient
	}

	return serviceClient, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Authorization", "Bearer "+c.token)

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
	if policy == nil {
		policy = func(resp *http.Response) bool {
			return resp.StatusCode == 429
		}
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

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
