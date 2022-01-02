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

resource "akcauth_authorization_code_client" "toto" {
  client_id = "toto"
  client_name = "toto name"
  allowed_scopes = ["toto.read", "toto.write"]
  redirect_uris = ["https://callback"]
}

resource "akcauth_api_scope" "basic_read" {
  name = "basic.read"
}