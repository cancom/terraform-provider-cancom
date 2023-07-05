package client_cmsmgw

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
)

// GetConnection - Returns a specific connection
func (c *Client) GetConnection(connectionID string) (*Connection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ipsec/%s", c.HostURL, connectionID), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	/*err2 := fmt.Errorf("status: %d, body: %s", 404, body)
	  if (err != nil) || (err2 != nil) {
	      return nil, err2
	  }*/
	if err != nil {
		return nil, err
	}

	connection := Connection{}
	err = json.Unmarshal(body, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

// CreateConnection - Create new connection
func (c *Client) CreateConnection(connection *ConnectionCreateRequest) (*Connection, error) {
	rb, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ipsec", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	new_connection := Connection{}
	err = json.Unmarshal(body, &new_connection)
	if err != nil {
		return nil, err
	}

	return &new_connection, nil
}

// UpdateConnection - update a connection
func (c *Client) UpdateConnection(connectionId string, connection *ConnectionUpdateRequest) (*Connection, error) {
	rb, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/ipsec/%s", c.HostURL, connectionId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	new_connection := Connection{}
	err = json.Unmarshal(body, &new_connection)
	if err != nil {
		return nil, err
	}

	return &new_connection, nil
}

// Delete Connection - delete connection
func (c *Client) DeleteConnection(connectionID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/ipsec/%s", c.HostURL, connectionID), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	new_connection := Connection{}
	err = json.Unmarshal(body, &new_connection)

	return nil
}
