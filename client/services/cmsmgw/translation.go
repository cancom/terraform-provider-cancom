package client_cmsmgw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
)

// Get Translation
func (c *Client) GetTranslation(translationId string) (*Translation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/translation/%s", c.HostURL, translationId), nil)
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	translation := Translation{}
	err = json.Unmarshal(body, &translation)
	if err != nil {
		return nil, err
	}

	return &translation, nil
}

// Create Translation
func (c *Client) CreateTranslation(translation *TranslationCreateRequest) (*Translation, error) {
	rb, err := json.Marshal(translation)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/translation", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	new_translation := Translation{}
	err = json.Unmarshal(body, &new_translation)
	if err != nil {
		return nil, err
	}

	return &new_translation, nil
}

// Update Translation
func (c *Client) UpdateTranslation(translationId string, translation *TranslationUpdateRequest) (*Translation, error) {
	rb, err := json.Marshal(translation)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/translation/%s", c.HostURL, translationId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return nil, err
	}

	new_translation := Translation{}
	err = json.Unmarshal(body, &new_translation)
	if err != nil {
		return nil, err
	}

	return &new_translation, nil
}

func (c *Client) DeleteTranslation(translationId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/translation/%s", c.HostURL, translationId), nil)
	if err != nil {
		return err
	}

	body, err := (*client.Client)(c).DoRequest(req)
	if err != nil {
		return err
	}

	new_translation := Translation{}
	err = json.Unmarshal(body, &new_translation)

	if (new_translation.DeploymentState != "deleted") && (err != nil) { // Status muss noch definiert werden
		return errors.New(string(body))
	}

	return nil
}
