package client

import (
	"encoding/json"
	"fmt"
)

func (c *Client) CreateApiResource(name string, model *ApiResourceCreate) error {
	req, err := c.prepareRequest("PUT", fmt.Sprintf("%s/%s/resources/%s", c.HostURL, c.Config.BasePath, name), model)
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
	req, err := c.prepareRequest("GET", fmt.Sprintf("%s/%s/resources/%s", c.HostURL, c.Config.BasePath, apiResourceName), nil)
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
	req, err := c.prepareRequest("POST", fmt.Sprintf("%s/%s/resources/%s", c.HostURL, c.Config.BasePath, apiResourceName), model)
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
	req, err := c.prepareRequest("DELETE", fmt.Sprintf("%s/%s/resources/%s", c.HostURL, c.Config.BasePath, apiResourceName), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
