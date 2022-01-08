package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ApiResource_Create_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(200, "", func(req *http.Request) {
		assert.Equal(t, "PUT", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources", AnyTestHostUrl), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"name":"api-resource-name","displayName":"display-name","scopes":["basic.read","basic.write"]}`, string(requestContent))
		callbacked = true
	})

	model := ApiResourceCreate{
		Name:        "api-resource-name",
		DisplayName: "display-name",
		Scopes:      []string{"basic.read", "basic.write"},
	}

	err := c.CreateApiResource(&model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_ApiResource_Get_ShouldPass(t *testing.T) {
	callbacked := false
	responseJson := `{"name":"api-resource-name","displayName":"display-name","scopes":["basic.read","basic.write"]}`
	c := setupWithCallback(200, responseJson, func(req *http.Request) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources/%s", AnyTestHostUrl, "api-resource-name"), req.URL.String())
		callbacked = true
	})

	model, err := c.GetApiResource("api-resource-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
	assert.Equal(t, "api-resource-name", model.Name)
	assert.Equal(t, "display-name", model.DisplayName)
	assert.Equal(t, []string{"basic.read", "basic.write"}, model.Scopes)
}

func Test_ApiResource_Delete_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources/%s", AnyTestHostUrl, "api-resource-name"), req.URL.String())
		callbacked = true
	})

	err := c.DeleteApiResource("api-resource-name")

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_ApiResource_UpdateAllProperties_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources/%s", AnyTestHostUrl, "api-resource-name"), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"displayName":"updated","scopes":["s1"]}`, string(requestContent))
		callbacked = true
	})

	model := ApiResourceUpdate{
		DisplayName: "updated",
		Scopes:      []string{"s1"},
	}

	err := c.UpdateApiResource("api-resource-name", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_ApiResource_Update_DisplayName_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources/%s", AnyTestHostUrl, "api-resource-name"), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"displayName":"updated"}`, string(requestContent))
		callbacked = true
	})

	model := ApiResourceUpdate{
		DisplayName: "updated",
	}

	err := c.UpdateApiResource("api-resource-name", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}

func Test_ApiResource_Update_Scopes_ShouldPass(t *testing.T) {
	callbacked := false
	c := setupWithCallback(204, "", func(req *http.Request) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/api/resources/%s", AnyTestHostUrl, "api-resource-name"), req.URL.String())
		requestContent, _ := ioutil.ReadAll(req.Body)
		assert.Equal(t, `{"scopes":["s1"]}`, string(requestContent))
		callbacked = true
	})

	model := ApiResourceUpdate{
		Scopes: []string{"s1"},
	}

	err := c.UpdateApiResource("api-resource-name", &model)

	assert.Nil(t, err)
	assert.True(t, callbacked)
}
