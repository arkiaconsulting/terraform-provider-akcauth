package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateApiResource(model *ApiResourceCreate) error {
	req, err := c.prepareRequest("PUT", fmt.Sprintf("%s/api/resources", c.HostURL), model)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetApiResource(apiResourceName string) (*ApiResource, error) {
	req, err := c.prepareRequest("GET", fmt.Sprintf("%s/api/resources/%s", c.HostURL, apiResourceName), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResource := ApiResource{}
	err = json.Unmarshal(body, &apiResource)
	if err != nil {
		return nil, err
	}

	return &apiResource, nil
}

func (c *Client) UpdateApiResource(apiResourceName string, model *ApiResourceUpdate) error {
	req, err := c.prepareRequest("POST", fmt.Sprintf("%s/api/resources/%s", c.HostURL, apiResourceName), model)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteApiResource(apiResourceName string) error {
	req, err := c.prepareRequest("DELETE", fmt.Sprintf("%s/api/resources/%s", c.HostURL, apiResourceName), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
