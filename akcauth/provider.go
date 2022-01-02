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
			"azuread_audience": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_AUTH_AUDIENCE", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"akcauth_authorization_code_client": resourceAuthorizationCodeClient(),
			"akcauth_api_scope":                 resourceApiScope(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	serverBaseUrl := d.Get("server_url").(string)
	audience := d.Get("azuread_audience").(string)

	var diags diag.Diagnostics

	if (serverBaseUrl == "") || (audience == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create client",
			Detail:   "One of the provider configuration settings is missing",
		})

		return nil, diags
	}

	config := client.ClientConfig{
		HostUrl:    serverBaseUrl,
		ResourceId: audience,
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
