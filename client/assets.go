package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const AnyTestHostUrl string = "https://any"

var (
	GetDoFunc        func(req *http.Request) (*http.Response, error)
	GetAuthorizeFunc func(req *http.Request, config *ClientConfig) error
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

type MockAuthorizer struct {
	AuthorizeFunc func(req *http.Request) error
}

type requestCallback func(*http.Request)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func (m *MockAuthorizer) Authorize(req *http.Request, config *ClientConfig) error {
	return GetAuthorizeFunc(req, config)
}

func jsonToBody(json string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(json)))
}

func setup(responseStatusCode int, responseJson string) *Client {
	HttpClient = &MockClient{}
	Authorizer = &MockAuthorizer{}
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: responseStatusCode,
			Body:       jsonToBody(responseJson),
		}, nil
	}
	GetAuthorizeFunc = func(req *http.Request, config *ClientConfig) error {
		return nil
	}
	config := ClientConfig{
		HostUrl:    AnyTestHostUrl,
		ResourceId: "any",
		BasePath:   "testpath",
	}

	c, _ := NewClient(&config)

	return c
}

func setupWithCallback(responseStatusCode int, responseJson string, callback requestCallback) *Client {
	HttpClient = &MockClient{}
	Authorizer = &MockAuthorizer{}
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		callback(req)
		return &http.Response{
			StatusCode: responseStatusCode,
			Body:       jsonToBody(responseJson),
		}, nil
	}
	GetAuthorizeFunc = func(req *http.Request, config *ClientConfig) error {
		return nil
	}
	config := ClientConfig{
		HostUrl:    AnyTestHostUrl,
		ResourceId: "any",
		BasePath:   "testpath",
	}

	c, _ := NewClient(&config)

	return c
}

func anyRequest() *http.Request {
	req, _ := http.NewRequest("any", "any", strings.NewReader("any"))

	return req
}
