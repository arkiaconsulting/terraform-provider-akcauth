package akcauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
			"allowed_grant_types": {
				Type:     schema.TypeList,
				Computed: true,
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
	clientId := d.Get("client_id").(string)
	log.Printf("[INFO] Creating resource (%s)", clientId)

	c := m.(*client.Client)

	var diags diag.Diagnostics

	existing, err := c.GetAuthorizationCodeClient(clientId)
	if err != nil {
		cErr, ok := err.(*client.ClientError)
		if !ok || (ok && cErr.Status != http.StatusNotFound) {
			return diag.FromErr(fmt.Errorf("checking for presence of existing %s: %+v", clientId, err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("A resource with the ID %q already exists - to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", existing.ClientId, "akcauth_authorization_code_client"))
	}

	clientName := d.Get("client_name").(string)
	allowedScopesRaw := d.Get("allowed_scopes").([]interface{})
	redirectUrisRaw := d.Get("redirect_uris").([]interface{})

	allowedScopes := expandString(allowedScopesRaw)
	redirectUris := expandString(redirectUrisRaw)

	allowedGrantTypes := [1]string{"client_credentials"}
	model := client.AuthorizationCodeClientCreate{
		ClientName:        clientName,
		AllowedScopes:     allowedScopes,
		RedirectUris:      redirectUris,
		AllowedGrantTypes: allowedGrantTypes[:],
	}

	err = c.CreateAuthorizationCodeClient(clientId, &model)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clientId)

	resourceAuthorizationCodeClientRead(ctx, d, m)

	return diags
}

func resourceAuthorizationCodeClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading resource (%s)", d.Id())

	c := m.(*client.Client)

	var diags diag.Diagnostics

	clientId := d.Id()

	authCodeClient, err := c.GetAuthorizationCodeClient(clientId)
	if err != nil {
		cErr, ok := err.(*client.ClientError)
		if ok {
			if cErr.Status == http.StatusNotFound {
				log.Printf("[WARN] The authorization code client was (%s) not found, removing from state", d.Id())
				d.SetId("")
				return nil
			} else {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(err)
		}
	}

	d.Set("client_id", authCodeClient.ClientId)
	d.Set("client_name", authCodeClient.ClientName)
	d.Set("allowed_scopes", authCodeClient.AllowedScopes)
	d.Set("redirect_uris", authCodeClient.RedirectUris)
	d.Set("enabled", authCodeClient.Enabled)
	d.Set("allowed_grant_types", authCodeClient.AllowedGrantTypes)

	return diags
}

func resourceAuthorizationCodeClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating resource (%s)", d.Id())

	c := m.(*client.Client)

	clientId := d.Id()

	authCodeClient, err := c.GetAuthorizationCodeClient(clientId)
	if err != nil {
		_, ok := err.(*client.ClientError)
		if ok {
			return diag.FromErr(err)
		} else {
			return diag.FromErr(err)
		}
	}

	// Could use json-patch instead
	updateModel := authCodeClient.ToUpdateModel()

	if d.HasChange("client_name") {
		updateModel.ClientName = d.Get("client_name").(string)
	}

	if d.HasChange("allowed_scopes") {
		allowedScopesRaw := d.Get("allowed_scopes").([]interface{})
		updateModel.AllowedScopes = expandString(allowedScopesRaw)
	}

	if d.HasChange("redirect_uris") {
		redirectUrisRaw := d.Get("redirect_uris").([]interface{})
		updateModel.RedirectUris = expandString(redirectUrisRaw)
	}

	if d.HasChange("enabled") {
		updateModel.Enabled = d.Get("enabled").(bool)
	}

	allowedGrantTypesRaw := d.Get("allowed_grant_types").([]interface{})
	updateModel.AllowedGrantTypes = expandString(allowedGrantTypesRaw)

	err = c.UpdateAuthorizationCodeClient(clientId, &updateModel)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAuthorizationCodeClientRead(ctx, d, m)
}

func resourceAuthorizationCodeClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting resource (%s)", d.Id())

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
