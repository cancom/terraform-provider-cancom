package client_dynamiccloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

var urlPath = "v1/Projects"

func (c *Client) GetVpcProject(id string) (*VpcProject, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		// this is required to allow us to use this function to detect successful deletion
		// we are polling describe until we get NotFound
		if strings.HasPrefix(err.Error(), "status: 404") {
			return nil, nil
		}
		return nil, err
	}

	vpcProject := VpcProject{}
	err = json.Unmarshal(body, &vpcProject)
	if err != nil {
		return nil, err
	}

	return &vpcProject, nil
}

func (c *Client) CreateVpcProject(vpcProjectCreateRequest *VpcProjectCreateRequest) (*VpcProject, error) {
	body, err := json.Marshal(vpcProjectCreateRequest)
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

	vpcProject := VpcProject{}
	err = json.Unmarshal(resp, &vpcProject)
	if err != nil {
		return nil, err
	}

	return &vpcProject, nil
}

func (c *Client) DeleteVpcProject(id string) error {
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
