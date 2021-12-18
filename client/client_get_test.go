package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCodeClient_Get_ShouldPass(t *testing.T) {
	callbacked := false
	responseJson := `{"clientId":"client-id","CLIENTNAME":"client-name","AllowedScopes":["s1","s2"],"redirectUris":["r1","r2"],"enabled":true,"requireClientSecret":true,"requirePkce":false,"allowedGrantTypes":["g1"],"allowOfflineAccess":true,"clientSecrets":[{"value":"s1","type":"SharedSecret"},{"value":"s2","type":"SharedSecret"}]}`
	c := setupWithCallback(200, responseJson, func(req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/clients/%s", AnyTestHostUrl, "client-id"), req.URL.String())
		callbacked = true
	})

	model, err := c.GetAuthorizationCodeClient("client-id")

	assert.Nil(t, err)
	assert.True(t, callbacked)
	assert.Equal(t, "client-id", model.ClientId)
	assert.Equal(t, "client-name", model.ClientName)
	assert.Equal(t, []string{"s1", "s2"}, model.AllowedScopes)
	assert.Equal(t, []string{"r1", "r2"}, model.RedirectUris)
	assert.True(t, model.Enabled)
	assert.True(t, model.RequireClientSecret)
	assert.False(t, model.RequirePkce)
	assert.Equal(t, []string{"g1"}, model.AllowedGrantTypes)
	assert.True(t, model.AllowOfflineAccess)
	assert.EqualValues(t, []ClientSecret{{Value: "s1", Type: "SharedSecret"}, {Value: "s2", Type: "SharedSecret"}}, model.ClientSecrets)
}
