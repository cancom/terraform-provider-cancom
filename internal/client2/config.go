package client2

import (
	"fmt"
	"log"
	"os"

	"github.com/cancom/terraform-provider-cancom/internal/sdk"
)

type ClientConfig struct {
	Token     *string
	Endpoints map[string]string
}

type Client struct {
	connections map[string]*sdk.CancomClient
	endpoints   map[string]string
	token       string
}

// Configures a new client implementation
// The token is retrieved in the following order, higher means that this token is being used
// 1. Token variable provided in the config
// 2. Token retrieved from the profile
// 3. Token retrieved from the
func NewClient(config *ClientConfig) (*Client, error) {
	// if token provided
	if config.Token != nil && *config.Token != "" {
		return &Client{}, nil
	}

	// Profile currently not supporeted

	// if token in environment variables
	token, exists := os.LookupEnv("CANCOM_ACCESS_TOKEN")
	if exists && token != "" {
		return &Client{}, nil
	}

	return nil, fmt.Errorf("a token is required")
}

func (c Client) GetClient(serviceName string) *sdk.CancomClient {
	log.Println("Test")
	mutexKey := serviceName
	GlobalMutexKV.Lock(mutexKey)
	defer GlobalMutexKV.Unlock(mutexKey)

	if connection, ok := c.connections[serviceName]; ok {
		return connection
	}

	service := SupportedServices[serviceName]

	endpoint := c.endpoints[serviceName]
	connection := service.Factory(c.token, &endpoint)

	c.connections[serviceName] = connection

	return connection
}
