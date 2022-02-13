package akcauth

import (
	"context"
	"terraform-provider-akcauth/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiScopeCreate,
		ReadContext:   resourceApiScopeRead,
		DeleteContext: resourceApiScopeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceApiScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	showInDiscoveryDocument := d.Get("show_in_discovery").(bool)
	userClaims := d.Get("user_claims").([]string)
	properties := d.Get("properties").(map[string]string)
	enabled := d.Get("enabled").(bool)
	required := d.Get("required").(bool)
	emphasize := d.Get("emphasize").(bool)

	model := client.ApiScopeCreate{
		DisplayName:             displayName,
		Description:             description,
		ShowInDiscoveryDocument: showInDiscoveryDocument,
		UserClaims:              userClaims,
		Properties:              properties,
		Enabled:                 enabled,
		Required:                required,
		Emphasize:               emphasize,
	}

	err := c.CreateApiScope(scopeName, &model)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(scopeName)

	resourceApiScopeRead(ctx, d, m)

	return diags
}

func resourceApiScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Id()

	apiScope, err := c.GetApiScope(scopeName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", apiScope.Name)

	return diags
}

func resourceApiScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Id()

	err := c.DeleteApiScope(scopeName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
