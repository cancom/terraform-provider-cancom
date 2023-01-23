package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HTTPMethod int64

var MainEndpoint = "portal.cancom.io"

const (
	GET HTTPMethod = iota
	POST
	PUT
	UPDATE
	DELETE
)

func (m HTTPMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	}

	return "unknown"
}

type HttpClient interface {
	Configure(HostURL string, Token string)
	DoRequest(req *http.Request) ([]byte, error)
}

type CancomClient struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

type Operation struct {
	Name   string
	Method HTTPMethod
	Path   string
}

func NewCancomClient(token, endpoint string) (*CancomClient, error) {
	return &CancomClient{
		Token:      token,
		HostURL:    endpoint,
		HTTPClient: &http.Client{},
	}, nil
}

func (c *CancomClient) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "C_HMAC_SHA256 "+c.Token)
	url, err := url.Parse(fmt.Sprintf("https://%s%s", c.HostURL, req.URL.Path))
	if err != nil {
		return nil, err
	}

	req.URL = url

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("[WARN] CancomClient - Connection could not be closed")
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *CancomClient) Post(path string, input interface{}, resultBuffer interface{}) (*interface{}, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &resultBuffer)
	if err != nil {
		return nil, err
	}

	return &resultBuffer, nil
}

func (c *CancomClient) Get(path string, outputBuffer interface{}) (*interface{}, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &outputBuffer)
	if err != nil {
		return nil, err
	}

	return &outputBuffer, nil
}

func (c *CancomClient) Put(path string, input interface{}, output interface{}) (*interface{}, error) {
	rb, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", path, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (c *CancomClient) Delete(path string, outputBuffer *interface{}) (*interface{}, error) {
	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, outputBuffer)
	if err != nil {
		return nil, err
	}

	return outputBuffer, nil
}
