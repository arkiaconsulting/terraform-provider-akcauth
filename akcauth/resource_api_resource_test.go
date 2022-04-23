package akcauth

import (
	"fmt"
	"log"
	"terraform-provider-akcauth/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ApiResourceResource struct{}

func TestAccApiResource_EnsureAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_resource", "basic_api")
	r := ApiResourceResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.basic_api"),
					resource.TestCheckResourceAttr("akcauth_api_resource.basic_api", "display_name", fmt.Sprintf("acctest apiresource %d", data.RandomInteger)),
				),
			},
		},
	})
}

func TestAccApiResource_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_resource", "basic_api")
	r := ApiResourceResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.basic_api"),
				),
			},
			{
				Config: r.displayNameUpdate(data, "new display name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("akcauth_api_resource.basic_api", "display_name", "new display name"),
				),
			},
		},
	})
}

func TestAccApiResource_CanBeImported(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_resource", "basic_api")
	r := ApiResourceResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.basic_api"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccResource_NoLongerExists(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_resource", "basic_api")
	r := ApiResourceResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiResourceResourceExist(t, "akcauth_api_resource.basic_api"),
					testAccCheckApiResourceDisappears("akcauth_api_resource.basic_api"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func (r ApiResourceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_api_resource" "basic_api" {
		name = "acctest-apiresource-%d"
		display_name = "acctest apiresource %d"
		scopes = [ "api_read_%d", "api_write_%d" ]
	}
	`, base(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiResourceResource) displayNameUpdate(data acceptance.TestData, displayName string) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_api_resource" "basic_api" {
		name = "acctest-apiresource-%d"
		display_name = "%s"
		scopes = [ "api_read_%d", "api_write_%d" ]
	}
	`, base(data), data.RandomInteger, displayName, data.RandomInteger, data.RandomInteger)
}

func testAccCheckApiResourceResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] Ensure that the Api resource (%s) exists", resourceName)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		apiResourceName := rs.Primary.Attributes["name"]

		_, err := getTestClient().GetApiResource(apiResourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckApiResourceDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		log.Print("[INFO] Ensure that all the Api scope resources were destroyed")

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_api_resource" {
				continue
			}

			apiResourceName := rs.Primary.ID
			_, err := getTestClient().GetApiResource(name)
			if err == nil {
				return fmt.Errorf("The Api resource with name '%s' still exists. (%s)", apiResourceName, name)
			}
		}

		return nil
	}
}

func testAccCheckApiResourceDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] Test is manually deleting the Api resource (%s)", resourceName)
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		apiResourceName := resourceState.Primary.ID

		err := getTestClient().DeleteApiResource(apiResourceName)
		if err != nil {
			return fmt.Errorf("We were unable to delete the remote Api resource '%s'", apiResourceName)
		}

		return nil
	}
}
