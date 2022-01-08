terraform {
  required_providers {
    akcauth = {
      version = "~>0.2"
      source  = "github.com/arkiaconsulting/akcauth"
    }
  }
}

provider "akcauth" {
  server_url = "https://auth.arkia.dev"
  azuread_audience = "api://arkia-identity"
}

resource "akcauth_authorization_code_client" "my_client" {
  client_id = "myClient"
  client_name = "My client"
  allowed_scopes = [akcauth_api_scope.api1_read.name, akcauth_api_scope.api1_write.name]
  redirect_uris = ["https://callback"]
}

resource "akcauth_api_scope" "api1_read" {
  name = "api1.read"
}

resource "akcauth_api_scope" "api1_write" {
  name = "api1.write"
}

resource "akcauth_api_resource" "api1" {
  name = "api1"
  display_name = "Api 1"
  scopes = [
    akcauth_api_scope.api1_read.name,
    akcauth_api_scope.api1_write.name,
  ]
}