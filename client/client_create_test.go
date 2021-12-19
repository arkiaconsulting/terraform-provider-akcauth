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
		assert.Equal(t, fmt.Sprintf("%s/api/clients", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"clientId":"client-id","clientName":"client name","allowedScopes":["basic","readwrite"],"redirectUris":["https://callback"]}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientCreate{
		ClientId:      "client-id",
		ClientName:    "client name",
		AllowedScopes: []string{"basic", "readwrite"},
		RedirectUris:  []string{"https://callback"},
	}

	err := c.CreateAuthorizationCodeClient(&model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
