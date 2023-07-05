package client_cmsmgw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
)

// ----------------------------------Gateway-----------------------------------
// GetGateway - Return a spcific gateway
func (c *Client) GetGateway(mgwID string) (*Gateway, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/mgw/%s", c.HostURL, mgwID), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)

	if err != nil {
		return nil, err
	}

	mgw := Gateway{}
	err = json.Unmarshal(body, &mgw)
	if err != nil {
		return nil, err
	}
	//err2 := fmt.Errorf("status: %d, body: %s", 404, body)
	//return nil, err2
	return &mgw, nil
}

// CreateMgw - Create new mgw
func (c *Client) CreateGateway(gateway *GatewayCreateRequest) (*Gateway, error) {
	rb, err := json.Marshal(gateway)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/mgw", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	mgw := Gateway{}
	err = json.Unmarshal(body, &mgw)
	if err != nil {
		return nil, err
	}

	return &mgw, nil
}

// UpdateMgw - update a mgw
func (c *Client) UpdateGateway(mgwID string, gateway *GatewayUpdateRequest) (*Gateway, error) {
	rb, err := json.Marshal(gateway)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/mgw/%s", c.HostURL, mgwID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	mgw := Gateway{}
	err = json.Unmarshal(body, &mgw)
	if err != nil {
		return nil, err
	}

	return &mgw, nil
}

// Delete mgw
func (c *Client) DeleteGateway(mgwID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/mgw/%s", c.HostURL, mgwID), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	mgw := Gateway{}
	err = json.Unmarshal(body, &mgw)

	if (mgw.State != "deleted") && (err != nil) { // Status muss noch definiert werden
		return errors.New(string(body))
	}

	return nil
}
