package akcauth

import (
	"context"
	"terraform-provider-akcauth/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuthorizationCodeClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthorizationCodeClientCreate,
		ReadContext:   resourceAuthorizationCodeClientRead,
		UpdateContext: resourceAuthorizationCodeClientUpdate,
		DeleteContext: resourceAuthorizationCodeClientDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAuthorizationCodeClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	clientId := d.Get("client_id").(string)
	clientName := d.Get("client_name").(string)
	allowedScopesRaw := d.Get("allowed_scopes").([]interface{})
	allowedScopes := make([]string, len(allowedScopesRaw))
	for i, raw := range allowedScopesRaw {
		allowedScopes[i] = raw.(string)
	}

	redirectUrisRaw := d.Get("redirect_uris").([]interface{})
	redirectUris := make([]string, len(redirectUrisRaw))
	for i, raw := range redirectUrisRaw {
		redirectUris[i] = raw.(string)
	}

	model := client.AuthorizationCodeClientCreate{
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

	resourceAuthorizationCodeClientRead(ctx, d, m)

	return diags
}

func resourceAuthorizationCodeClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	clientId := d.Id()

	authCodeClient, err := c.GetAuthorizationCodeClient(clientId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("client_id", authCodeClient.ClientId)
	d.Set("client_name", authCodeClient.ClientName)
	d.Set("allowed_scopes", authCodeClient.AllowedScopes)
	d.Set("redirect_uris", authCodeClient.RedirectUris)
	d.Set("enabled", authCodeClient.Enabled)

	return diags
}

func resourceAuthorizationCodeClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	clientId := d.Id()

	updateModel := client.AuthorizationCodeClientUpdate{}

	if d.HasChange("client_name") {
		updateModel.ClientName = d.Get("client_name").(string)
	}

	if d.HasChange("allowed_scopes") {
		allowedScopesRaw := d.Get("allowed_scopes").([]interface{})
		allowedScopes := make([]string, len(allowedScopesRaw))
		for i, raw := range allowedScopesRaw {
			allowedScopes[i] = raw.(string)
		}
		updateModel.AllowedScopes = allowedScopes
	}

	if d.HasChange("redirect_uris") {
		redirectUrisRaw := d.Get("redirect_uris").([]interface{})
		redirectUris := make([]string, len(redirectUrisRaw))
		for i, raw := range redirectUrisRaw {
			redirectUris[i] = raw.(string)
		}
		updateModel.RedirectUris = redirectUris
	}

	if d.HasChange("enabled") {
		updateModel.Enabled = d.Get("enabled").(bool)
	}

	err := c.UpdateAuthorizationCodeClient(clientId, &updateModel)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAuthorizationCodeClientRead(ctx, d, m)
}

func resourceAuthorizationCodeClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	clientId := d.Id()

	err := c.DeleteAuthorizationCodeClient(clientId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
