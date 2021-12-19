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

func TestAccAuthorizationCodeClient_EnsureAttributes(t *testing.T) {
	data := randomAuthorizationCodeClient()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizationCodeClientResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "client_id", data.ClientId),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "client_name", data.ClientName),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "allowed_scopes.0", data.AllowedScopes[0]),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "allowed_scopes.1", data.AllowedScopes[1]),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "redirect_uris.0", data.RedirectUris[0]),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "redirect_uris.1", data.RedirectUris[1]),
				),
			},
		},
	})
}

func TestAccAuthorizationCodeClient_EnabledByDefault(t *testing.T) {
	data := randomAuthorizationCodeClient()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizationCodeClientResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
					resource.TestCheckResourceAttr("akcauth_authorization_code_client.my_client", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAuthorizationCodeClient_CanBeImported(t *testing.T) {
	data := randomAuthorizationCodeClient()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthorizationCodeClientResourceDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAuthorizationCodeClientResource_full(data),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthorizationCodeClientResourceExist(t, "akcauth_authorization_code_client.my_client"),
				),
			},
			{
				ResourceName:      "akcauth_authorization_code_client.my_client",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAuthorizationCodeClientResource_full(data *AuthorizationCodeClientTestData) string {
	return fmt.Sprintf(`
provider "akcauth" {}

resource "akcauth_authorization_code_client" "my_client" {
	client_id = "%s"
	client_name = "%s"
	allowed_scopes = [ "%s" ]
	redirect_uris = [ "%s" ]
}
`, data.ClientId, data.ClientName, strings.Join(data.AllowedScopes, `","`), strings.Join(data.RedirectUris, `","`))
}

func testAccCheckAuthorizationCodeClientResourceExist(t *testing.T, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		clientId := rs.Primary.Attributes["client_id"]

		_, err := c.GetAuthorizationCodeClient(clientId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckAuthorizationCodeClientResourceDestroy(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		c := testAccProvider.Meta().(*client.Client)

		for name, rs := range s.RootModule().Resources {
			if rs.Type != "akcauth_authorization_code_client" {
				continue
			}

			clientId := rs.Primary.ID
			_, err := c.GetAuthorizationCodeClient(clientId)
			if err == nil {
				return fmt.Errorf("The authorization code client with client Id '%s' still exists. (%s)", clientId, name)
			}
		}

		return nil
	}
}

type AuthorizationCodeClientTestData struct {
	ClientId      string
	ClientName    string
	AllowedScopes []string
	RedirectUris  []string
}

func randomAuthorizationCodeClient() *AuthorizationCodeClientTestData {
	data := AuthorizationCodeClientTestData{
		ClientId:      acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		ClientName:    acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
		AllowedScopes: make([]string, 2),
		RedirectUris:  make([]string, 2),
	}

	data.AllowedScopes[0] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	data.AllowedScopes[1] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	data.RedirectUris[0] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	data.RedirectUris[1] = acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	return &data
}
