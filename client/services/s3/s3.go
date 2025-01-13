package client_s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

// ----------------------- Bucket -------------------
func (c *Client) GetBucket(bucketName string) (*Bucket, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Buckets/%s", c.HostURL, bucketName), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	record := Bucket{}
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *Client) CreateBucket(bucket *Bucket) (*Bucket, error) {
	body, err := json.Marshal(bucket)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Buckets", c.HostURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	recordBody := Bucket{}
	err = json.Unmarshal(resp, &recordBody)
	if err != nil {
		return nil, err
	}

	return &recordBody, nil
}

func (c *Client) DeleteBucket(bucket string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Buckets/%s", c.HostURL, bucket), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// ----------------------- Users -------------------
func (c *Client) GetUser(userId string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, userId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	record := User{}
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (c *Client) CreateUser(user *UserCreateRequest) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Users", c.HostURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	recordBody := User{}
	err = json.Unmarshal(resp, &recordBody)
	if err != nil {
		return nil, err
	}

	return &recordBody, nil
}

func (c *Client) DeleteUser(userId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, userId), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateUser(userId string, request *UserUpdateRequest) (*User, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, userId), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = (*client.Client)(c).DoRequestWithRetry(req, nil)
	if err != nil {
		return nil, err
	}

	userResponse := User{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse, nil
}
