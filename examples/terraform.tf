terraform {
  required_providers {
    akcauth = {
      version = "~>0.2"
      source  = "github.com/arkiaconsulting/akcauth"
    }
  }
}

provider "akcauth" {
  server_url         = "https://akc-duende.azurewebsites.net/"
  api_base_path      = "/my"
  authorization_type = "client_credentials"
  client_id          = "client"
  client_secret      = "secret"
  scopes             = ["IdentityServerApi"]
  # azuread_audience = "api://arkia-identity"
}

resource "akcauth_authorization_code_client" "my_client" {
  client_id      = "myClient"
  client_name    = "My client"
  allowed_scopes = [akcauth_api_scope.my_resource_read.name, akcauth_api_scope.my_resource_write.name]
  redirect_uris  = ["https://callback"]
}

resource "akcauth_api_scope" "my_resource_read" {
  name = "myResource.read"
}

resource "akcauth_api_scope" "my_resource_write" {
  name = "myResource.write"
}

resource "akcauth_api_resource" "my_resource" {
  name         = "myResource"
  display_name = "My resource"
  scopes = [
    akcauth_api_scope.my_resource_read.name,
    akcauth_api_scope.my_resource_write.name,
  ]
}
