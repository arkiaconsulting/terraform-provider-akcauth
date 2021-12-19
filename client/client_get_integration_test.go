package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_ShouldPass(t *testing.T) {
	config := ClientConfig{
		HostUrl:    os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		ResourceId: os.Getenv("AKC_AUTH_AUDIENCE"),
	}
	clientId := "toto"
	HttpClient = &http.Client{Timeout: 10 * time.Second}
	c, err := NewClient(&config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	err = c.CreateAuthorizationCodeClient(&AuthorizationCodeClientCreate{
		ClientId:   clientId,
		ClientName: "client name",
	})
	if err != nil {
		t.Logf("Client creation resulted in '%s'", err.Error())
	}

	myClient, err := c.GetAuthorizationCodeClient(clientId)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, clientId, myClient.ClientId)

	newClientName, _ := uuid.GenerateUUID()
	err = c.UpdateAuthorizationCodeClient(clientId, &AuthorizationCodeClientUpdate{ClientName: newClientName})
	if err != nil {
		t.Logf("Client update resulted in '%s'", err.Error())
	}

	err = c.DeleteAuthorizationCodeClient(clientId)
	assert.Nil(t, err)
}
