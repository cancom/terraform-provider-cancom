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
	role        string
	mu          *sync.Mutex
	tp          *tokenProvider
}

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	tp         *tokenProvider
}

func NewClient(host, token *string, role string) (*CcpClient, error) {
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

	origProvider := &tokenProvider{
		mu:    &sync.Mutex{},
		token: *token,
	}

	srClient := newHttpClient(serviceRegistry, origProvider)

	tp := &tokenProvider{
		mu:    &sync.Mutex{},
		token: *token,
		role:  role,
	}

	c := CcpClient{
		client:      &http.Client{Timeout: 2 * time.Minute},
		services:    map[string]*Client{"service-registry": srClient},
		serviceURLs: map[string]string{},
		mu:          &sync.Mutex{},
		initialized: false,
		role:        role,
		token:       *token,
		tp:          tp,
	}

	// if role is set, initialize the client and assume role
	if role != "" {
		c.initialize()
		iamClient, err := c.getSingleService("iam")
		if err != nil {
			return nil, err
		}

		tp.iamClient = newHttpClient(iamClient.ServiceEndpoint.Backend, origProvider)
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

func newHttpClient(host string, tp *tokenProvider) *Client {
	return &Client{
		HostURL:    host,
		tp:         tp,
		HTTPClient: &http.Client{Timeout: 2 * time.Minute},
	}
}

func (c *CcpClient) WithToken(token string) *CcpClient {
	serviceMap := map[string]*Client{}

	for name, service := range c.services {
		serviceMap[name] = newHttpClient(service.HostURL, c.tp)
	}

	return &CcpClient{
		client:      c.client,
		serviceURLs: c.serviceURLs,
		token:       token,
		mu:          c.mu,
		initialized: c.initialized,
		services:    serviceMap,
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

		serviceClient = newHttpClient(serviceUrl, c.tp)
		c.services[name] = serviceClient
	}

	return serviceClient, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	token, err := c.tp.GetToken()
	if err != nil {
		return nil, err
	}

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
	if policy == nil {
		policy = func(resp *http.Response) bool {
			return resp.StatusCode == 429
		}
	}

	token, err := c.tp.GetToken()
	if err != nil {
		return nil, err
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

		token, err := c.tp.GetToken()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
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
