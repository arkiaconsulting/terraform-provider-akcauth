package akcauth

import (
	"context"
	"terraform-provider-akcauth/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_BASE_ADDRESS", nil),
			},
			"api_base_path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_BASE_PATH", "/"),
			},
			"azuread_audience": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_AUDIENCE", nil),
			},
			"authorization_type": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_AUTHORIZATION_TYPE", "client_credentials"),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_CLIENT_SECRET", nil),
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"akcauth_authorization_code_client": resourceAuthorizationCodeClient(),
			"akcauth_api_scope":                 resourceApiScope(),
			"akcauth_api_resource":              resourceApiResource(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure(),
	}
}

func providerConfigure() func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		serverBaseUrl := d.Get("server_url").(string)
		apiBasePath := d.Get("api_base_path").(string)
		audience := d.Get("azuread_audience").(string)
		authorizationType := d.Get("authorization_type").(string)
		clientId := d.Get("client_id").(string)
		clientSecret := d.Get("client_secret").(string)
		scopesRaw := d.Get("scopes").([]interface{})

		scopes := make([]string, len(scopesRaw))
		for i, raw := range scopesRaw {
			scopes[i] = raw.(string)
		}

		var diags diag.Diagnostics

		if serverBaseUrl == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create client",
				Detail:   "One of the provider configuration settings is missing",
			})

			return nil, diags
		}

		config := client.ClientConfig{
			HostUrl:           serverBaseUrl,
			ResourceId:        audience,
			AuthorizationType: authorizationType,
			ClientId:          clientId,
			ClientSecret:      clientSecret,
			Scopes:            scopes,
			BasePath:          apiBasePath,
		}
		c, err := client.NewClient(&config)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create client",
				Detail:   "Unable to instanciate the client using the given configuration",
			})

			return nil, diags
		}

		return c, diags
	}
}
