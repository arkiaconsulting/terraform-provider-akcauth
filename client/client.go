package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type ClientConfig struct {
	HostUrl    string
	ResourceId string
}

type Client struct {
	HostURL    string
	HTTPClient HttpClientWrapper
	ResourceId string
	Token      string
}

type HttpClientWrapper interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient HttpClientWrapper
)

func init() {
	HttpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// NewClient -
func NewClient(config *ClientConfig) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("client configuration is missing")
	}

	if config.HostUrl == "" {
		return nil, fmt.Errorf("the host url is missing from the client configuration")
	}

	c := Client{
		HTTPClient: HttpClient,
		HostURL:    config.HostUrl,
		ResourceId: config.ResourceId,
	}

	return &c, nil
}

func authorizeRequest(req *http.Request, c *Client) error {
	if c.Token == "" {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return fmt.Errorf("cannot get credential: %+v", err)
		}

		token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{c.ResourceId},
		})
		if err != nil {
			return fmt.Errorf("cannot get Azure access token: %+v", err)
		}
		c.Token = token.Token
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	err := authorizeRequest(req, c)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) CreateAuthorizationCodeClient(model *AuthorizationCodeClientCreate) error {
	rb, err := json.Marshal(model)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/clients", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateAuthorizationCodeClient(clientId string, model *AuthorizationCodeClientUpdate) error {
	rb, err := json.Marshal(model)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/clients/%s", c.HostURL, clientId), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteAuthorizationCodeClient(clientId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/clients/%s", c.HostURL, clientId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetAuthorizationCodeClient(clientId string) (*AuthorizationCodeClient, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/clients/%s", c.HostURL, clientId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	acClient := AuthorizationCodeClient{}
	err = json.Unmarshal(body, &acClient)
	if err != nil {
		return nil, err
	}

	return &acClient, nil
}
