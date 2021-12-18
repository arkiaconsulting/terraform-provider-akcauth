package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCodeClient_UpdateAllProperties_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"clientName":"client name","allowedScopes":["basic","readwrite"],"redirectUris":["https://callback"],"enabled":true}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientUpdate{
		ClientName:    "client name",
		AllowedScopes: []string{"basic", "readwrite"},
		RedirectUris:  []string{"https://callback"},
		Enabled:       true,
	}

	err := c.UpdateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_AuthorizationCodeClient_UpdateClientName_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"clientName":"client name"}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientUpdate{
		ClientName: "client name",
	}

	err := c.UpdateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_AuthorizationCodeClient_UpdateAllowedScopes_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"allowedScopes":["basic","readwrite"]}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientUpdate{
		AllowedScopes: []string{"basic", "readwrite"},
	}

	err := c.UpdateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_AuthorizationCodeClient_UpdateRedirectUris_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"redirectUris":["https://callback","https://callback2"]}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientUpdate{
		RedirectUris: []string{"https://callback", "https://callback2"},
	}

	err := c.UpdateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_AuthorizationCodeClient_UpdateEnabled_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/client-id", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"enabled":true}`, string(requestContent))
		callbacked = true
	})

	model := AuthorizationCodeClientUpdate{
		Enabled: true,
	}

	err := c.UpdateAuthorizationCodeClient("client-id", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
