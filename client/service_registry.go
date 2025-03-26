package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *CcpClient) getAllServices() ([]Service, error) {
	serviceRegistry, ok := c.services["service-registry"]
	if !ok {
		return nil, fmt.Errorf("service-registry not initialized")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Services", serviceRegistry.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := serviceRegistry.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var services []Service
	err = json.Unmarshal(body, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (c *CcpClient) getSingleService(name string) (*Service, error) {
	serviceRegistry, ok := c.services["service-registry"]
	if !ok {
		return nil, fmt.Errorf("service-registry not initialized")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Services/%s", serviceRegistry.HostURL, name), nil)
	if err != nil {
		return nil, err
	}

	body, err := serviceRegistry.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var service *Service
	err = json.Unmarshal(body, &service)
	if err != nil {
		return nil, err
	}

	return service, nil
}
