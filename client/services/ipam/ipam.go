package client_ipam

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

// GetHost - Returns a specifc host address
func (c *Client) GetHost(hostId string) (*Host, error) {
	//return nil, errors.New(c.HostURL)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Hosts/%s", c.HostURL, hostId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	host := Host{}
	err = json.Unmarshal(body, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

// CreateHost - Create new host
func (c *Client) CreateHost(host *HostCreateRequest) (*Host, error) {
	rb, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	// sleep for a random time to resolve race condition problem
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Second)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Hosts", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	newHost := Host{}
	err = json.Unmarshal(body, &newHost)
	if err != nil {
		return nil, err
	}

	return &newHost, nil
}

// UpdateHost - Updates an Host
func (c *Client) UpdateHost(hostId string, host *HostUpdateRequest) (*Host, error) {
	rb, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Hosts/%s", c.HostURL, hostId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	newHost := Host{}
	err = json.Unmarshal(body, &newHost)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newHost, nil
}

// DeleteHost - Release a host
func (c *Client) DeleteHost(hostId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Hosts/%s", c.HostURL, hostId), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	if string(body) != "Deleted host" {
		return nil //errors.New(string(body))
	}

	return nil
}
