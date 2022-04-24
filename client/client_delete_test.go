package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthorizationCodeClient_Delete_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/my/clients/%s", AnyTestHostUrl, "client-id"), req.URL.String())
		callbacked = true
	})

	err := c.DeleteAuthorizationCodeClient("client-id")

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
