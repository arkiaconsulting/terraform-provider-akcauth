package akcauth

import (
	"fmt"
	"terraform-provider-akcauth/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApiScope_EnsureAttributes(t *testing.T) {
	data := randomApiScope()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiScopeDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApiScopeResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiScopeResourceExist(t, "akcauth_api_scope.basic_read"),
					resource.TestCheckResourceAttr("akcauth_api_scope.basic_read", "name", data.Name),
				),
			},
		},
	})
}

func TestAccApiScope_CanBeImported(t *testing.T) {
	data := randomApiScope()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiScopeDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApiScopeResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiScopeResourceExist(t, "akcauth_api_scope.basic_read"),
				),
			},
			{
				ResourceName:      "akcauth_api_scope.basic_read",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiScopeResource_full(data *ApiScopeTestData) string {
	return fmt.Sprintf(`
provider "akcauth" {}

resource "akcauth_api_scope" "basic_read" {
	name = "%s"
}
`, data.Name)
}

func testAccCheckApiScopeResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		scopeName := rs.Primary.Attributes["name"]

		_, err := c.GetApiScope(scopeName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckApiScopeDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_api_scope" {
				continue
			}

			scopeName := rs.Primary.ID
			_, err := c.GetApiScope(scopeName)
			if err == nil {
				return fmt.Errorf("The Api scope with name '%s' still exists. (%s)", scopeName, name)
			}
		}

		return nil
	}
}

type ApiScopeTestData struct {
	Name string
}

func randomApiScope() *ApiScopeTestData {
	data := ApiScopeTestData{
		Name: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
	}

	return &data
}
