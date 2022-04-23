package akcauth

import (
	"fmt"
	"log"
	"terraform-provider-akcauth/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ApiScopeResource struct{}

func TestAccApiScope_EnsureAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_scope", "basic_read")
	r := ApiScopeResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiScopeDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiScopeResourceExist(t, "akcauth_api_scope.basic_read"),
				),
			},
		},
	})
}

func TestAccApiScope_CanBeImported(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_scope", "basic_read")
	r := ApiScopeResource{}

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

func TestAccApiScope_NoLongerExists(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_scope", "basic_read")
	r := ApiScopeResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiScopeDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiScopeResourceExist(t, "akcauth_api_scope.basic_read"),
					testAccCheckApiScopeDisappears("akcauth_api_scope.basic_read"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccApiScope_RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_api_scope", "basic_read")
	r := ApiScopeResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckApiScopeDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiScopeResourceExist(t, "akcauth_api_scope.basic_read"),
				),
			},
			{
				Config:      r.requiresImport(data),
				ExpectError: RequiresImportError("akcauth_api_scope"),
			},
		},
	})
}

func (r ApiScopeResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_api_scope" "basic_read" {
		name = "acctest-apiscope-%d"
	}
	`, base(data), data.RandomInteger)
}

func (r ApiScopeResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "akcauth_api_scope" "import" {
		name = akcauth_api_scope.basic_read.name
	}
	`, r.basic(data))
}

func testAccCheckApiScopeResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] Ensure that the Api scope (%s) exists", resourceName)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		scopeName := rs.Primary.Attributes["name"]

		_, err := getTestClient().GetApiScope(scopeName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckApiScopeDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		log.Print("[INFO] Ensure that all the Api scope resources were destroyed")

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_api_scope" {
				continue
			}

			scopeName := rs.Primary.ID
			_, err := getTestClient().GetApiScope(scopeName)
			if err == nil {
				return fmt.Errorf("The Api scope with name '%s' still exists. (%s)", scopeName, name)
			}
		}

		return nil
	}
}

func testAccCheckApiScopeDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] Test is manually deleting the Api scope (%s)", resourceName)
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		apiScopeName := resourceState.Primary.ID

		err := getTestClient().DeleteApiScope(apiScopeName)
		if err != nil {
			return fmt.Errorf("We were unable to delete the remote Api scope '%s'", apiScopeName)
		}

		return nil
	}
}
