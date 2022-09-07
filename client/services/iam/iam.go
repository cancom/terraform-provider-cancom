package client_iam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
)

type Client client.Client

func (c *Client) GetUser(id string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) CreateUser(userCreateRequest *UserCreateRequest) (*UserCreateResponse, error) {
	body, err := json.Marshal(userCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Users", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	user := UserCreateResponse{}
	err = json.Unmarshal(resp, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) UpdateUser(id string, userUpdateRequest *UserUpdateRequest) error {
	body, err := json.Marshal(userUpdateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteUser(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Users/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetServiceUser(id string) (*ServiceUser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/ServiceUsers/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	serviceUser := ServiceUser{}
	err = json.Unmarshal(body, &serviceUser)
	if err != nil {
		return nil, err
	}

	return &serviceUser, nil
}

func (c *Client) CreateServiceUser(serviceUserCreateRequest *ServiceUserCreateRequest) (*ServiceUserCreateResponse, error) {
	body, err := json.Marshal(serviceUserCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/ServiceUsers", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	serviceUser := ServiceUserCreateResponse{}
	err = json.Unmarshal(resp, &serviceUser)
	if err != nil {
		return nil, err
	}

	return &serviceUser, nil
}

func (c *Client) UpdateServiceUser(id string, serviceUserUpdateRequest *ServiceUserUpdateRequest) error {
	body, err := json.Marshal(serviceUserUpdateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/ServiceUsers/%s", c.HostURL, id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteServiceUser(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/ServiceUsers/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRole(id string) (*Role, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Roles/%s", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	role := Role{}
	err = json.Unmarshal(body, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (c *Client) CreateRole(roleCreateRequest *RoleCreateRequest) (*RoleCreateResponse, error) {
	body, err := json.Marshal(roleCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Roles", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	role := RoleCreateResponse{}
	err = json.Unmarshal(resp, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (c *Client) UpdateRole(id string, roleUpdateRequest *RoleUpdateRequest) error {
	body, err := json.Marshal(roleUpdateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Roles/%s", c.HostURL, id), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteRole(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Roles/%s", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AssignPolicyToUser(policyRequest *PolicyRequest, principal string) error {
	body, err := json.Marshal(policyRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/PolicyDocuments/%s", c.HostURL, principal), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil

}

func (c *Client) RemovePolicyFromUser(PolicyRequest *PolicyRequest, principal string) error {
	body, err := json.Marshal(PolicyRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/PolicyDocuments/%s?services=%s", c.HostURL, principal, PolicyRequest.Service), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
