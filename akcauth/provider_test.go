package akcauth

import (
	"log"
	"os"
	"strings"
	"terraform-provider-akcauth/acceptance"
	"terraform-provider-akcauth/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider = Provider()
var testAccProviders = testAccProvidersFactory()

func testAccProvidersFactory() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"akcauth": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func init() {
	log.Print("[INFO] Initializing the test provider")
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func getTestClient() *client.Client {
	log.Print("[INFO] Creating test client")

	scopes := make([]string, 1)
	scopes[0] = "IdentityServerApi"

	config := client.ClientConfig{
		HostUrl:           os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		BasePath:          os.Getenv("AKC_AUTH_BASE_PATH"),
		AuthorizationType: os.Getenv("AKC_AUTH_AUTHORIZATION_TYPE"),
		ClientId:          os.Getenv("AKC_AUTH_CLIENT_ID"),
		ClientSecret:      os.Getenv("AKC_AUTH_CLIENT_SECRET"),
		Scopes:            strings.Split(" ", os.Getenv("AKC_AUTH_SCOPES")),
	}

	c, err := client.NewClient(&config)
	if err != nil {
		panic(err.Error())
	}

	return c
}

func testAccPreCheck(t *testing.T) {
	log.Print("[INFO] Checking test pre-requisites")

	if v := os.Getenv("AKC_AUTH_BASE_ADDRESS"); v == "" {
		t.Fatal("the AKC_AUTH_BASE_ADDRESS environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_BASE_PATH"); v == "" {
		t.Fatal("the AKC_AUTH_BASE_PATH environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_AUTHORIZATION_TYPE"); v == "" {
		t.Fatal("the AKC_AUTH_AUTHORIZATION_TYPE environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_CLIENT_ID"); v == "" {
		t.Fatal("the AKC_AUTH_CLIENT_ID environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_CLIENT_SECRET"); v == "" {
		t.Fatal("the AKC_AUTH_CLIENT_SECRET environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_SCOPES"); v == "" {
		t.Fatal("the AKC_AUTH_SCOPES environment variable must be set for acceptance tests")
	}
}

func base(data acceptance.TestData) string {
	return `
provider "akcauth" {
	api_base_path = "/my"
	authorization_type = "client_credentials"
	client_id = "client"
	client_secret = "secret"
	scopes = [ "IdentityServerApi" ]
}
`
}
