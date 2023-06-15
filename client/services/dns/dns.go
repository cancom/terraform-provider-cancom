package client_dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) GetZone(id string) (*Zone, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Zones/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	zone := Zone{}
	err = json.Unmarshal(body, &zone)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

func (c *Client) GetAllZones() ([]Zone, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Zones", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	var zones []Zone
	err = json.Unmarshal(body, &zones)
	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (c *Client) GetRecord(id string, zoneName string) (*Record, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Records/%s?zoneName=%s", c.HostURL, id, zoneName), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	record := Record{}
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *Client) GetAllRecords() ([]Record, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Records", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	var records []Record
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *Client) CreateRecord(record *RecordCreateRequest) (*Record, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Records", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	recordBody := Record{}
	err = json.Unmarshal(body, &recordBody)
	if err != nil {
		return nil, err
	}

	return &recordBody, nil
}

func (c *Client) UpdateRecord(id string, record *RecordUpdateRequest) (*Record, error) {
	body, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Records/%s", c.HostURL, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	recordBody := Record{}
	err = json.Unmarshal(body, &recordBody)
	if err != nil {
		return nil, err
	}

	return &recordBody, nil
}

func (c *Client) DeleteRecord(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Records/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
