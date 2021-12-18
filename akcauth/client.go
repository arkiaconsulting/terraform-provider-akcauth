package akcauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ClientConfig struct {
	HostUrl string
}

type Client struct {
	HostURL    string
	HTTPClient HttpClientWrapper
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
		Token:      "todo",
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)

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
