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
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func jsonToBody(json string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(json)))
}

func setup(responseStatusCode int, responseJson string) *Client {
	HttpClient = &MockClient{}
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: responseStatusCode,
			Body:       jsonToBody(responseJson),
		}, nil
	}
	config := ClientConfig{
		HostUrl: AnyTestHostUrl,
	}

	c, _ := NewClient(&config)

	return c
}

func anyRequest() *http.Request {
	req, _ := http.NewRequest("any", "any", strings.NewReader("any"))

	return req
}
