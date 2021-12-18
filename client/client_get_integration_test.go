package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Integration_Get_ShouldPass(t *testing.T) {
	config := ClientConfig{
		HostUrl:    os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		ResourceId: os.Getenv("AKC_AUTH_RESOURCE_ID"),
	}
	HttpClient = &http.Client{Timeout: 10 * time.Second}
	c, err := NewClient(&config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	_, err = c.GetAuthorizationCodeClient("testing")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	myClient, err := c.GetAuthorizationCodeClient("testing")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, "testing", myClient.ClientId)
}
