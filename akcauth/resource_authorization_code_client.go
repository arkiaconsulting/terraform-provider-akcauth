package akcauth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuthorizationCodeClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorizationCodeClientCreate,
		ReadContext:   resourceAuthorizationCodeClientRead,
		UpdateContext: resourceAuthorizationCodeClientUpdate,
		DeleteContext: resourceAuthorizationCodeClientDelete,
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceAuthorizationCodeClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	c := m.(Client)

	var diags diag.Diagnostics

	clientId := d.Get("client_id").(string)
	clientName := d.Get("client_name").(string)
	allowedScopes := d.Get("allowed_scopes").([]string)
	redirectUris := d.Get("redirect_uris").([]string)

	model := AuthorizationCodeClientCreate{
		ClientId:      clientId,
		ClientName:    clientName,
		AllowedScopes: allowedScopes,
		RedirectUris:  redirectUris,
	}

	err := c.CreateAuthorizationCodeClient(&model)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clientId)

	return diags
}

func resourceAuthorizationCodeClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceAuthorizationCodeClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAuthorizationCodeClientRead(ctx, d, m)
}

func resourceAuthorizationCodeClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}