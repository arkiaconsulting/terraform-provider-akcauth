package akcauth

import (
	"fmt"
	"log"
	"terraform-provider-akcauth/acceptance"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type ClientResource struct{}

func TestAccAuthorizationCodeClient_EnsureAttributes(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_authorization_code_client", "my_client")
	r := ClientResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAuthorizationCodeClient_Update_ClientName(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_authorization_code_client", "my_client")
	r := ClientResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
				),
			},
			{
				Config: r.nameUpdate(data, "new-client-name"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "client_name", "new-client-name"),
				),
			},
		},
	})
}

func TestAccAuthorizationCodeClient_Update_Disable(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_authorization_code_client", "my_client")
	r := ClientResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
				),
			},
			{
				Config: r.disable(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccAuthorizationCodeClient_NoLongerExists(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_authorization_code_client", "my_client")
	r := ClientResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
					testAccCheckAuthorizationCodeClientDisappears("akcauth_authorization_code_client.my_client"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAuthorizationCodeClient_RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "akcauth_authorization_code_client", "my_client")
	r := ClientResource{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
				),
			},
			{
				Config:      r.requiresImport(data),
				ExpectError: RequiresImportError("akcauth_authorization_code_client"),
			},
		},
	})
}

func testAccCheckAuthorizationCodeClientResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		log.Printf("[INFO] Ensure that the authorization code client (%s) exists", rs.Primary.ID)

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		clientId := rs.Primary.ID

		_, err := getTestClient().GetAuthorizationCodeClient(clientId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAuthorizationCodeClientResourceDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		log.Print("[INFO] Ensure that all the client resources were destroyed")

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_authorization_code_client" {
				continue
			}

			clientId := rs.Primary.ID
			_, err := getTestClient().GetAuthorizationCodeClient(clientId)
			if err == nil {
				return fmt.Errorf("The authorization code client with client Id '%s' still exists. (%s)", clientId, name)
			}
		}

		return nil
	}
}

func (r ClientResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_authorization_code_client" "my_client" {
		client_id = "client-%d"
		client_name = "name-%d"
		allowed_scopes = [ "scope-%d" ]
		redirect_uris = [ "https://host/callback-%d" ]
		enabled = true
	}
	`, base(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ClientResource) nameUpdate(data acceptance.TestData, name string) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_authorization_code_client" "my_client" {
		client_id = "client-%d"
		client_name = "%s"
		allowed_scopes = [ "scope-%d" ]
		redirect_uris = [ "https://host/callback-%d" ]
		enabled = true
	}
	`, base(data), data.RandomInteger, name, data.RandomInteger, data.RandomInteger)
}

func (r ClientResource) disable(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s
	
	resource "akcauth_authorization_code_client" "my_client" {
		client_id = "client-%d"
		client_name = "name-%d"
		allowed_scopes = [ "scope-%d" ]
		redirect_uris = [ "https://host/callback-%d" ]
		enabled = false
	}
	`, base(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ClientResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

	resource "akcauth_authorization_code_client" "import" {
		client_id = akcauth_authorization_code_client.my_client.client_id
		client_name = akcauth_authorization_code_client.my_client.client_name
		allowed_scopes = akcauth_authorization_code_client.my_client.allowed_scopes
		redirect_uris = akcauth_authorization_code_client.my_client.redirect_uris
		enabled = akcauth_authorization_code_client.my_client.enabled
	}
	`, r.basic(data))
}

func testAccCheckAuthorizationCodeClientDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] Test is manually deleting the authorization code client (%s)", resourceName)
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		clientId := resourceState.Primary.ID

		err := getTestClient().DeleteAuthorizationCodeClient(clientId)
		if err != nil {
			return fmt.Errorf("We were unable to delete the remote authorization code client '%s'", clientId)
		}

		return nil
	}
}
