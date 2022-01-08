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
	c := setupWithCallback(200, "", func(req *http.Request) {
		assert.Equal(t, "PUT", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/scopes", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"name":"scope-name"}`, string(requestContent))
		callbacked = true
	})

	model := ApiScopeCreate{
		Name: "scope-name",
	}

	err := c.CreateApiScope(&model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_Scope_Get_ShouldPass(t *testing.T) {
	callbacked := false
	responseJson := `{"name":"scope-name"}`
	c := setupWithCallback(200, responseJson, func(req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/scopes/%s", AnyTestHostUrl, "scope-name"), req.URL.String())
		callbacked = true
	})

	model, err := c.GetApiScope("scope-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
	assert.Equal(t, "scope-name", model.Name)
}

func Test_Scope_Delete_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/scopes/%s", AnyTestHostUrl, "scope-name"), req.URL.String())
		callbacked = true
	})

	err := c.DeleteApiScope("scope-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
