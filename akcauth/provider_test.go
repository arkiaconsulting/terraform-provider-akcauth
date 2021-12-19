package akcauth

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"akcauth": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AKC_AUTH_BASE_ADDRESS"); v == "" {
		t.Fatal("the AKC_AUTH_BASE_ADDRESS environment variable must be set for acceptance tests")
	}

	if v := os.Getenv("AKC_AUTH_AUDIENCE"); v == "" {
		t.Fatal("the AKC_AUTH_AUDIENCE environment variable must be set for acceptance tests")
	}
}
