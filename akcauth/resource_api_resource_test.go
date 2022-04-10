package akcauth

import (
	"fmt"
	"strings"
	"terraform-provider-akcauth/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApiResource_EnsureAttributes(t *testing.T) {
	data := randomApiResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApiResourceResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.my_resource"),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "name", data.Name),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "display_name", data.DisplayName),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "scopes.0", data.Scopes[0]),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "scopes.1", data.Scopes[1]),
				),
			},
		},
	})
}

func TestAccApiResource_Update(t *testing.T) {
	data := randomApiResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApiResourceResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.my_resource"),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "name", data.Name),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "display_name", data.DisplayName),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "scopes.0", data.Scopes[0]),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "scopes.1", data.Scopes[1]),
				),
			},
			{
				Config: testAccApiResourceResource_single(data.Name, "updated", []string{"s1"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "name", data.Name),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "display_name", "updated"),
					resource.TestCheckResourceAttr("akcauth_api_resource.my_resource", "scopes.0", "s1"),
				),
			},
		},
	})
}

func TestAccApiResource_CanBeImported(t *testing.T) {
	data := randomApiResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccApiResourceResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.my_resource"),
				),
			},
			{
				ResourceName:      "akcauth_api_resource.my_resource",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiResourceResource_full(data *ApiResourceTestData) string {
	return fmt.Sprintf(`
provider "akcauth" {}
%s
`, testAccApiResourceResource_single(data.Name, data.DisplayName, data.Scopes))
}

func testAccApiResourceResource_single(name string, displayName string, scopes []string) string {
	return fmt.Sprintf(`
resource "akcauth_api_resource" "my_resource" {
	name = "%s"
	display_name = "%s"
	scopes = [ "%s" ]
}
`, name, displayName, strings.Join(scopes, `","`))
}

func testAccCheckApiResourceResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		apiResourceName := rs.Primary.Attributes["name"]

		_, err := c.GetApiResource(apiResourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckApiResourceResourceDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_api_resource" {
				continue
			}

			apiResourceName := rs.Primary.ID
			_, err := c.GetApiResource(apiResourceName)
			if err == nil {
				return fmt.Errorf("The api resource with name '%s' still exists. (%s)", apiResourceName, name)
			}
		}

		return nil
	}
}

type ApiResourceTestData struct {
	Name        string
	DisplayName string
	Scopes      []string
}

func randomApiResource() *ApiResourceTestData {
	data := ApiResourceTestData{
		Name:        acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		DisplayName: acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		Scopes:      make([]string, 2),
	}

	data.Scopes[0] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	data.Scopes[1] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	return &data
}
