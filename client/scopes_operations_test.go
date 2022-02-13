package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Scope_Create_ShouldPass(t *testing.T) {
	callbacked := false
	scopeName := "my-scope"
	c := setupWithCallback(201, "", func(req *http.Request) {
		assert.Equal(t, "PUT", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/testpath/scopes/%s", AnyTestHostUrl, scopeName), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"displayName":"display-name","description":"description","showInDiscoveryDocument":true,"userClaims":["given_name"],"properties":{"prop":"value"},"enabled":true,"required":true,"emphasize":false}`, string(requestContent))
		callbacked = true
	})

	model := ApiScopeCreate{
		DisplayName:             "display-name",
		Description:             "description",
		ShowInDiscoveryDocument: true,
		UserClaims:              []string{"given_name"},
		Properties:              map[string]string{"prop": "value"},
		Enabled:                 true,
		Required:                true,
		Emphasize:               false,
	}

	err := c.CreateApiScope(scopeName, &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_Scope_Get_ShouldPass(t *testing.T) {
	callbacked := false
	responseJson := `{"name":"scope-name","displayName":"display-name","description":"description","showInDiscoveryDocument":true,"userClaims":["given_name"],"properties":{"prop":"value"},"enabled":true,"required":true,"emphasize":false}`
	c := setupWithCallback(200, responseJson, func(req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/testpath/scopes/%s", AnyTestHostUrl, "scope-name"), req.URL.String())
		callbacked = true
	})

	model, err := c.GetApiScope("scope-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
	assert.Equal(t, "scope-name", model.Name)
	assert.Equal(t, "display-name", model.DisplayName)
	assert.Equal(t, "description", model.Description)
	assert.Equal(t, true, model.ShowInDiscoveryDocument)
	assert.Equal(t, []string{"given_name"}, model.UserClaims)
	assert.Equal(t, map[string]string{"prop": "value"}, model.Properties)
	assert.Equal(t, true, model.Enabled)
	assert.Equal(t, true, model.Required)
	assert.Equal(t, false, model.Emphasize)
}

func Test_Scope_Delete_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/testpath/scopes/%s", AnyTestHostUrl, "scope-name"), req.URL.String())
		callbacked = true
	})

	err := c.DeleteApiScope("scope-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
