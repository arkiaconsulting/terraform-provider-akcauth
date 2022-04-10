package client

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Integration_ShouldPass(t *testing.T) {
	config := ClientConfig{
		HostUrl:           os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		BasePath:          os.Getenv("AKC_AUTH_BASE_PATH"),
		AuthorizationType: os.Getenv("AKC_AUTH_AUTHORIZATION_TYPE"),
		ClientId:          os.Getenv("AKC_AUTH_CLIENT_ID"),
		ClientSecret:      os.Getenv("AKC_AUTH_CLIENT_SECRET"),
		Scopes:            strings.Split(os.Getenv("AKC_AUTH_SCOPES"), " "),
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
		t.Fatalf("Client creation resulted in '%s'", err.Error())
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
