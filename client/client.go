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

	"golang.org/x/oauth2/clientcredentials"
)

type ClientConfig struct {
	HostUrl           string
	ResourceId        string
	ClientId          string
	ClientSecret      string
	AuthorizationType string
	Scopes            []string
	BasePath          string
}

type Client struct {
	HostURL    string
	HTTPClient HttpClientWrapper
	Authorizer RequestAuthorizer
	Config     *ClientConfig
}

type TokenProvider struct {
	Token string
}

type RequestAuthorizer interface {
	Authorize(req *http.Request, cfg *ClientConfig) error
}

type HttpClientWrapper interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient HttpClientWrapper
	Authorizer RequestAuthorizer
)

func init() {
	HttpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	Authorizer = &TokenProvider{}
}

// Authorize using Azure
func (tokenProvider *TokenProvider) Authorize(req *http.Request, cfg *ClientConfig) error {
	if tokenProvider.Token == "" {
		if cfg.AuthorizationType == "azure" {
			cred, err := azidentity.NewDefaultAzureCredential(nil)
			if err != nil {
				return fmt.Errorf("cannot get credential: %+v", err)
			}

			token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: []string{cfg.ResourceId},
			})
			if err != nil {
				return fmt.Errorf("cannot get Azure access token: %+v", err)
			}
			tokenProvider.Token = token.Token
		} else if cfg.AuthorizationType == "client_credentials" {
			conf := clientcredentials.Config{
				ClientID:     cfg.ClientId,
				ClientSecret: cfg.ClientSecret,
				TokenURL:     fmt.Sprintf("%s/connect/token", cfg.HostUrl),
				Scopes:       cfg.Scopes,
			}

			token, err := conf.Token(context.Background())
			if err != nil {
				return fmt.Errorf("cannot get access token using client credentials: %+v", err)
			}

			tokenProvider.Token = token.AccessToken
		} else {
			return fmt.Errorf("cannot get access token using authorization type: %+v", cfg.AuthorizationType)
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenProvider.Token))

	return nil
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
		Authorizer: Authorizer,
		HostURL:    config.HostUrl,
		Config:     config,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

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

func (c *Client) prepareRequest(method string, url string, model interface{}) (*http.Request, error) {
	rb, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	err = c.Authorizer.Authorize(req, c.Config)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) CreateAuthorizationCodeClient(model *AuthorizationCodeClientCreate) error {
	req, err := c.prepareRequest("PUT", fmt.Sprintf("%s/%s/clients", c.HostURL, c.Config.BasePath), model)
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
	req, err := c.prepareRequest("POST", fmt.Sprintf("%s/%s/clients/%s", c.HostURL, c.Config.BasePath, clientId), model)
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
	req, err := c.prepareRequest("DELETE", fmt.Sprintf("%s/%s/clients/%s", c.HostURL, c.Config.BasePath, clientId), nil)
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
	req, err := c.prepareRequest("GET", fmt.Sprintf("%s/%s/clients/%s", c.HostURL, c.Config.BasePath, clientId), nil)
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
