package client_cmsmgw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) CreateGateway(record *GatewayCreateRequest) (*Gateway, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/Records", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	recordBody := Gateway{}
	err = json.Unmarshal(body, &recordBody)
	if err != nil {
		return nil, err
	}

	return &recordBody, nil
}
