package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCodeClient_Create_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(200, "", func(req *http.Request) {
		assert.Equal(t, "PUT", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/my/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"clientName":"client name","allowedScopes":["basic","readwrite"],"redirectUris":["https://callback"],"allowedGrantTypes":["client_credentials"]}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientCreate{
		ClientName:        "client name",
		AllowedScopes:     []string{"basic", "readwrite"},
		RedirectUris:      []string{"https://callback"},
		AllowedGrantTypes: []string{"client_credentials"},
	}

	err := c.CreateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
