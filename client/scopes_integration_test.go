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

	err = c.CreateApiScope(&ApiScopeCreate{
		Name: scopeName,
	})
	if err != nil {
		t.Logf("Scope creation resulted in '%s'", err.Error())
	}

	myScope, err := c.GetApiScope(scopeName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, scopeName, myScope.Name)

	err = c.DeleteApiScope(scopeName)
	assert.Nil(t, err)
}
