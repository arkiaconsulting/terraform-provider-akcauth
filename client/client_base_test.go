package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewAkcAuthClientWithoutConfig(t *testing.T) {
	HttpClient = &MockClient{}

	_, err := NewClient(nil)

	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func Test_NewAkcAuthClientWithoutHostUrl(t *testing.T) {
	config := ClientConfig{}

	_, err := NewClient(&config)

	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func Test_NewAkcAuthClient(t *testing.T) {
	config := ClientConfig{
		HostUrl: AnyTestHostUrl,
	}

	_, err := NewClient(&config)

	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
}

func Test_Request_BadRequest_ShouldFail(t *testing.T) {
	c := setup(400, `{"errorMessage":"bad request message"}`)

	_, err := c.doRequest(anyRequest())

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("status: 400, body: %s", `{"errorMessage":"bad request message"}`), err.Error())
}

func Test_Request_Conflict_ShouldFail(t *testing.T) {
	c := setup(409, `{"errorMessage":"there is a conflict"}`)

	_, err := c.doRequest(anyRequest())

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("status: 409, body: %s", `{"errorMessage":"there is a conflict"}`), err.Error())
}

func Test_Request_NotFound_ShouldFail(t *testing.T) {
	c := setup(404, "")

	_, err := c.doRequest(anyRequest())

	assert.NotNil(t, err)
	assert.EqualValues(t, fmt.Sprintf("status: 404, body: %s", ""), err.Error())
}

func Test_Request_NoContent_ShouldPass(t *testing.T) {
	c := setup(204, "")

	_, err := c.doRequest(anyRequest())

	assert.Nil(t, err)
}

func Test_Request_Created_ShouldPass(t *testing.T) {
	c := setup(201, "")

	_, err := c.doRequest(anyRequest())

	assert.Nil(t, err)
}

func Test_Request_Ok_ShouldPass(t *testing.T) {
	c := setup(200, "content")

	_, err := c.doRequest(anyRequest())

	assert.Nil(t, err)
}
