package client_windowsos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"os"

	"slices"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

var urlPath = "api/v1/Deployment"

func (c *Client) GetWindowsDeployment(id string) (*WindowsOS_Deplyoment, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	windowosDeployment := WindowsOS_Deplyoment{}
	err = json.Unmarshal(body, &windowosDeployment)
	if err != nil {
		return nil, err
	}

	return &windowosDeployment, nil
}

func (c *Client) CreateWindowsDeployment(windowsOSDeployment *WindowsOS_Deplyoment) (*WindowsOS_Deplyoment, error) {
	body, err := json.Marshal(windowsOSDeployment)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.HostURL, urlPath), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	windowosDeployment := WindowsOS_Deplyoment{}
	err = json.Unmarshal(resp, &windowosDeployment)
	if err != nil {
		return nil, err
	}

	return &windowosDeployment, nil
}

// function to wait for succeded deployment
func (c *Client) CreateWindowsDeploymentStatus(id string) (*WindowsOS_Deplyoment, error) {

	errorstatus := []int{6, 5}
	sucessstatus := []int{4}

	timeoutCount := 0

	//validation ressource. Waits until the deployment of the software has been finished.
	//The deployment itself is run by the CANCOM Windows OS Service backend.
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)

		req.Header.Add("Content-Type", "application/json")

		resp, err := (*client.Client)(c).DoRequest(req)
		if err != nil {

			// allow timeouts because of long running queries in background
			if os.IsTimeout(err) {
				timeoutCount++
			} else {
				return nil, err
			}
			if timeoutCount > 10 {
				return nil, err
			}
		} else {
			timeoutCount = 0
			apiResultObject := WindowsOS_Deplyoment{}
			err = json.Unmarshal(resp, &apiResultObject)
			if err != nil {
				return nil, err
			}
			if slices.Contains(errorstatus, apiResultObject.Status) {
				return nil, err
			} else if slices.Contains(sucessstatus, apiResultObject.Status) {
				return &apiResultObject, nil
			}
			time.Sleep(30 * time.Second) // sleep for 30 seconds to aviod active waiting
		}

	}

}

func (c *Client) UpdateWindowsOsDeployment(id string, windowsOSDeployment *WindowsOS_Deplyoment) (*WindowsOS_Deplyoment, error) {
	body, err := json.Marshal(windowsOSDeployment)
	if err != nil {
		return nil, err
	}

	currentWindowsOSDeployment, getterError := c.GetWindowsDeployment(windowsOSDeployment.Id)
	if getterError != nil {
		return nil, err
	}

	if currentWindowsOSDeployment.CoustomerEnvironmentID != windowsOSDeployment.CoustomerEnvironmentID {
		return nil, err
	}
	if currentWindowsOSDeployment.Computer.Role != windowsOSDeployment.Computer.Role {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	osDeployment := WindowsOS_Deplyoment{}
	err = json.Unmarshal(resp, &osDeployment)
	if err != nil {
		return nil, err
	}

	return &osDeployment, nil
}

func (c *Client) DeleteWindowsDeployment(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", c.HostURL, urlPath, id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
