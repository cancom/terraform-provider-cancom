package client_iam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cancom/terraform-provider-cancom/client"
	"github.com/golang-jwt/jwt/v5"
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

func (c *Client) AssumeRole(role *AssumeRoleRequest) (*AssumeRoleResponse, error) {
	body, err := json.Marshal(role)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/AssumeRole", c.HostURL), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	assumeRole := AssumeRoleResponse{}
	err = json.Unmarshal(resp, &assumeRole)
	if err != nil {
		return nil, err
	}

	return &assumeRole, nil
}

func ParseSessionToken(tokenString string) (*SessionClaims, error) {
	var claims SessionClaims
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}

	return &claims, nil
}

func (c *Client) CreateSession(sessionCreateRequest *SessionCreateRequest) (*SessionCreateResponse, error) {
	body, err := json.Marshal(sessionCreateRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/Sessions", c.HostURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	session := SessionCreateResponse{}
	err = json.Unmarshal(resp, &session)
	if err != nil {
		return nil, err
	}

	if session.Jwt != "" {
		claims, err := ParseSessionToken(session.Jwt)
		if err == nil && claims != nil {
			session.Claims = claims
			session.SessionID = claims.SessionID
			session.PrincipalCRN = claims.PrincipalCRN
			if claims.ExpiresAt != nil {
				session.ExpiresAt = claims.ExpiresAt.Unix()
			}
		}
	}

	return &session, nil
}

func (c *Client) GetSession(principalCRN string, sessionID string) (*Session, error) {
	sessions, err := c.ListSessions(principalCRN)
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.SessionID == sessionID {
			return &session, nil
		}
	}

	return nil, &client.HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    []byte(fmt.Sprintf("session %s not found for principal %s", sessionID, principalCRN)),
	}
}

func (c *Client) ListSessions(principalCRN string) ([]Session, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/Sessions/%s", c.HostURL, principalCRN), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	var sessions []Session
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (c *Client) GetSessions(principalCRN string) ([]Session, error) {
	return c.ListSessions(principalCRN)
}

func (c *Client) UpdateSession(principalCRN string, sessionID string, sessionUpdateRequest *SessionUpdateRequest) error {
	body, err := json.Marshal(sessionUpdateRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/Sessions/%s/%s", c.HostURL, principalCRN, sessionID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSession(principalCRN string, sessionID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/Sessions/%s/%s", c.HostURL, principalCRN, sessionID), nil)
	if err != nil {
		return err
	}

	_, err = (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}
