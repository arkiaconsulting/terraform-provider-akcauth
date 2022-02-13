package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Scopes_Integration_ShouldPass(t *testing.T) {
	config := ClientConfig{
		HostUrl:    os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		ResourceId: os.Getenv("AKC_AUTH_AUDIENCE"),
	}
	scopeName := "basic.read"
	HttpClient = &http.Client{Timeout: 10 * time.Second}
	c, err := NewClient(&config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	apiScope := ApiScopeCreate{
		DisplayName:             "display-name",
		Description:             "description",
		ShowInDiscoveryDocument: true,
		UserClaims:              []string{"given_name"},
		Properties:              map[string]string{"prop": "value"},
		Enabled:                 true,
		Required:                true,
		Emphasize:               false,
	}

	err = c.CreateApiScope(scopeName, &apiScope)
	if err != nil {
		t.Logf("Scope creation resulted in '%s'", err.Error())
	}

	myScope, err := c.GetApiScope(scopeName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, scopeName, myScope.Name)
	assert.Equal(t, apiScope.DisplayName, myScope.DisplayName)
	assert.Equal(t, apiScope.Description, myScope.Description)
	assert.Equal(t, apiScope.ShowInDiscoveryDocument, myScope.ShowInDiscoveryDocument)
	assert.Equal(t, apiScope.UserClaims, myScope.UserClaims)
	assert.Equal(t, apiScope.Properties, myScope.Properties)
	assert.Equal(t, apiScope.Enabled, myScope.Enabled)
	assert.Equal(t, apiScope.Required, myScope.Required)
	assert.Equal(t, apiScope.Emphasize, myScope.Emphasize)

	err = c.DeleteApiScope(scopeName)
	assert.Nil(t, err)
}
