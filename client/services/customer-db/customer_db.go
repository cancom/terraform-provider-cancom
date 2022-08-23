package client_customerdb

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) GetTranslation(id, sourceServiceName, targetSerivceName string) (*Translation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Translations/%s?sourceServiceName=%s&targetServiceName=%s", c.HostURL, id, sourceServiceName, targetSerivceName), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	translation := Translation{}
	err = json.Unmarshal(body, &translation)
	if err != nil {
		return nil, err
	}

	return &translation, nil
}

func (c *Client) GetCustomer(id string) (*Customer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/Customers/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	customer := Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (c *Client) GetAllCustomers() ([]Customer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/CustomersOS", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	var customers []Customer
	err = json.Unmarshal(body, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}
