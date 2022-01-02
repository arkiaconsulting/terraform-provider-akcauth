package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateApiScope(model *ApiScopeCreate) error {
	req, err := c.prepareRequest("PUT", fmt.Sprintf("%s/api/scopes", c.HostURL), model)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetApiScope(scopeName string) (*ApiScope, error) {
	req, err := c.prepareRequest("GET", fmt.Sprintf("%s/api/scopes/%s", c.HostURL, scopeName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	acClient := ApiScope{}
	err = json.Unmarshal(body, &acClient)
	if err != nil {
		return nil, err
	}

	return &acClient, nil
}

func (c *Client) DeleteApiScope(scopeName string) error {
	req, err := c.prepareRequest("DELETE", fmt.Sprintf("%s/api/scopes/%s", c.HostURL, scopeName), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
