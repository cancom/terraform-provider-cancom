package client_sslmonitoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

var urlPath = "v1/SslMonitors"

func (c *Client) GetSslMonitor(id string) (*SslMonitor, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	sslMonitor := SslMonitor{}
	err = json.Unmarshal(body, &sslMonitor)
	if err != nil {
		return nil, err
	}

	return &sslMonitor, nil
}

func (c *Client) StartSslScan(id string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/startScan", c.HostURL, urlPath, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateSslMonitor(sslMonitorCreateRequest *SslMonitorCreateRequest) (*SslMonitor, error) {
	body, err := json.Marshal(sslMonitorCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, urlPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	sslMonitor := SslMonitor{}
	err = json.Unmarshal(resp, &sslMonitor)
	if err != nil {
		return nil, err
	}

	return &sslMonitor, nil
}

func (c *Client) UpdateSslMonitor(id string, sslMonitorCreateRequest *SslMonitorCreateRequest) (*SslMonitor, error) {
	body, err := json.Marshal(sslMonitorCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	sslMonitor := SslMonitor{}
	err = json.Unmarshal(resp, &sslMonitor)
	if err != nil {
		return nil, err
	}

	return &sslMonitor, nil
}

func (c *Client) DeleteSslMonitor(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
