package client_serviceregistry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) GetService(id string) (*Service, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Services/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	service := Service{}
	err = json.Unmarshal(body, &service)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (c *Client) GetAllServices() ([]Service, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Services", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
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

func (c *Client) CreateOrUpdateService(service *CreateServiceBody) error {
	body, err := json.Marshal(service)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/Services", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteService(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/Services/%s", c.HostURL, name), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
