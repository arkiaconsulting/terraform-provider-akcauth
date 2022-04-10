package akcauth

import (
	"log"
	"os"
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
	scopes := make([]string, 1)
	scopes[0] = "IdentityServerApi"

	config := client.ClientConfig{
		HostUrl:           os.Getenv("AKC_AUTH_BASE_ADDRESS"),
		BasePath:          "/my",
		ResourceId:        os.Getenv("AKC_AUTH_AUDIENCE"),
		AuthorizationType: "client_credentials",
		ClientId:          "client",
		ClientSecret:      "secret",
		Scopes:            scopes,
	}

	c, err := client.NewClient(&config)
	if err != nil {
		panic(err.Error())
	}

	return c
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AKC_AUTH_BASE_ADDRESS"); v == "" {
		t.Fatal("the AKC_AUTH_BASE_ADDRESS environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_AUDIENCE"); v == "" {
		t.Fatal("the AKC_AUTH_AUDIENCE environment variable must be set for acceptance tests")
	}
}
