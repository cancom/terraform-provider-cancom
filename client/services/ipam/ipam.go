package client_ipam

import (
	//"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

// -----------------------------Network---------------------------------

// GetNetwork - Returns a specifc Network
func (c *Client) GetNetwork(networkId string) (*Network, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Networks/%s", c.HostURL, networkId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	network := Network{}
	err = json.Unmarshal(body, &network)
	if err != nil {
		return nil, err
	}

	return &network, nil
}

// CreateNetwork - Assign a new network from a given supernet or supernetPool
func (c *Client) CreateNetwork(network *NetworkCreateRequest) (*Network, error) {
	rb, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Second)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Networks/assign", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return nil, err
	}

	newNetwork := Network{}
	err = json.Unmarshal(body, &newNetwork)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newNetwork, nil

}

// UpdateNetwork - Update a network
func (c *Client) UpdateNetwork(networkId string, network *NetworkUpdateRequest) (*Network, error) {
	rb, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Networks/%s", c.HostURL, networkId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	newNetwork := Network{}
	err = json.Unmarshal(body, &newNetwork)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newNetwork, nil
}

// DeleteNetwork - Deleting a network releases it to the pool
func (c *Client) DeleteNetwork(networkId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Networks/release/%s", c.HostURL, networkId), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return err
	}

	delResponse := NetworkDeleteResponse{}
	err = json.Unmarshal(body, &delResponse)
	if err != nil {
		return err
	}

	//if delResponse.Message != "released successfully" {
	//	return errors.New("message: " + delResponse.Message)
	//}

	return nil
}

// -----------------------------Supernet---------------------------------

// GetSupernet - Returns a specifc Supernet
func (c *Client) GetSupernet(supernetId string) (*Supernet, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Networks/%s", c.HostURL, supernetId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	supernet := Supernet{}
	err = json.Unmarshal(body, &supernet)
	if err != nil {
		return nil, err
	}

	return &supernet, nil
}

// CreateSupernet - Create new supernet
func (c *Client) CreateSupernet(supernet *SupernetCreateRequest) (*Supernet, error) {
	rb, err := json.Marshal(supernet)
	if err != nil {
		return nil, err
	}

	//return nil, errors.New("message create supernet")
	// sleep for a random time to resolve potential race condition problem
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Second)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Networks", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return nil, err
	}

	newSupernet := Supernet{}
	err = json.Unmarshal(body, &newSupernet)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newSupernet, nil

}

// UpdateSupernet - Updates a supernet
func (c *Client) UpdateSupernet(supernetId string, supernet *SupernetUpdateRequest) (*Supernet, error) {
	rb, err := json.Marshal(supernet)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Networks/%s", c.HostURL, supernetId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return nil, err
	}

	newSupernet := Supernet{}
	err = json.Unmarshal(body, &newSupernet)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newSupernet, nil
}

// DeleteSupernet - Delete a Supernet
func (c *Client) DeleteSupernet(supernetId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Networks/%s", c.HostURL, supernetId), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return err
	}

	delResponse := SupernetDeleteResponse{}
	err = json.Unmarshal(body, &delResponse)
	if err != nil {
		return err
	}

	//if delResponse.Message != "released successfully" {
	//	return errors.New("message: " + delResponse.Message)
	//}

	return nil
}

// -----------------------------Instance---------------------------------

// GetInstance - Returns a specifc Instance
func (c *Client) GetInstance(instanceId string) (*Instance, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Instances/%s", c.HostURL, instanceId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	instance := Instance{}
	err = json.Unmarshal(body, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// CreateInstance - Create new instance
func (c *Client) CreateInstance(instance *InstanceCreateRequest) (*Instance, error) {
	rb, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}
	// sleep for a random time to resolve potential race condition problem
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Second)

	//return nil, errors.New("message: " + c.HostURL)  //for debugging

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Instances", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	newInstance := Instance{}
	err = json.Unmarshal(body, &newInstance)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newInstance, nil

}

// UpdateInstance - Updates an instance
func (c *Client) UpdateInstance(instanceId string, instance *InstanceUpdateRequest) (*Instance, error) {
	rb, err := json.Marshal(instance)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Instances/%s", c.HostURL, instanceId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	newInstance := Instance{}
	err = json.Unmarshal(body, &newInstance)
	if err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)
	return &newInstance, nil
}

// DeleteInstance - Delete an Instance
func (c *Client) DeleteInstance(instanceId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Instances/%s", c.HostURL, instanceId), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	delResponse := InstanceDeleteResponse{}
	err = json.Unmarshal(body, &delResponse)
	if err != nil {
		return err
	}

	//if delResponse.Message != "released successfully" {
	//	return errors.New("message: " + delResponse.Message)
	//}

	return nil
}

// -------------------------------------Host--------------------------------------------------

// GetHost - Returns a specifc host address
func (c *Client) GetHost(hostId string) (*Host, error) {
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

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
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

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
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

	body, err := (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return err
	}

	delResponse := HostDeleteResponse{}
	err = json.Unmarshal(body, &delResponse)
	if err != nil {
		return err
	}
	// included to find a sporadic problem leading to resources not released in the backend
	if delResponse.Message != "Item released successfully" && delResponse.Message != "released successfully" {
		return errors.New("message: " + delResponse.Message)
	}

	return nil
}
