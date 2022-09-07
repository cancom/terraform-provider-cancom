package client_carrental

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) GetOrder(orderID string) (*CarOrder, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/CarOrders/%s", c.HostURL, orderID), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)

	if err != nil {
		return nil, err
	}

	order := CarOrder{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c Client) CreateOrder(orderItems CarOrderCreateRequest) (*CarOrder, error) {
	rb, err := json.Marshal(orderItems)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/CarOrders", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(&c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	order := CarOrder{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c Client) UpdateOrder(orderItems CarOrderUpdateRequest) (*CarOrder, error) {
	rb, err := json.Marshal(orderItems)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/CarOrders", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(&c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	order := CarOrder{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *Client) DeleteOrder(orderID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/CarOrders/%s", c.HostURL, orderID), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
