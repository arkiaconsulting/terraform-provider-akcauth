package akcauth

import (
	"context"
	"log"
	"net/http"
	"terraform-provider-akcauth/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiScopeCreate,
		ReadContext:   resourceApiScopeRead,
		DeleteContext: resourceApiScopeDelete,
		UpdateContext: resourceApiScopeUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"show_in_discovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"user_claims": {
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
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"emphasize": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Default:  map[string]string{},
			},
		},
	}
}

func resourceApiScopeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating Api scope (%s)", d.Get("name"))

	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	showInDiscoveryDocument := d.Get("show_in_discovery").(bool)
	userClaimsRaw := d.Get("user_claims").([]interface{})
	propertiesRaw := d.Get("properties").(map[string]interface{})
	enabled := d.Get("enabled").(bool)
	required := d.Get("required").(bool)
	emphasize := d.Get("emphasize").(bool)

	userClaims := make([]string, len(userClaimsRaw))
	for i, raw := range userClaimsRaw {
		userClaims[i] = raw.(string)
	}

	properties := map[string]string{}
	for key, element := range propertiesRaw {
		properties[key] = element.(string)
	}

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
	log.Printf("[INFO] Reading Api scope (%s)", d.Id())

	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Id()

	apiScope, err := c.GetApiScope(scopeName)
	if err != nil {
		log.Printf("[WARN] Api scope (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", apiScope.Name)
	d.Set("display_name", apiScope.DisplayName)
	d.Set("description", apiScope.Description)
	d.Set("show_in_discovery", apiScope.ShowInDiscoveryDocument)
	d.Set("enabled", apiScope.Enabled)
	d.Set("required", apiScope.Required)
	d.Set("emphasize", apiScope.Emphasize)

	return diags
}

func resourceApiScopeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Api scope (%s)", d.Id())

	var diags diag.Diagnostics

	return diags
}

func resourceApiScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Api scope (%s)", d.Id())

	c := m.(*client.Client)

	var diags diag.Diagnostics

	scopeName := d.Id()

	err := c.DeleteApiScope(scopeName)
	if err != nil {
		cErr, ok := err.(*client.ClientError)
		if ok {
			if cErr.Status == http.StatusBadRequest {
				log.Printf("[WARN] The Api scope (%s) could not be deleted (Error %d)", d.Id(), cErr.Status)
				return nil
			} else {
				return diag.FromErr(err)
			}
		} else {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}
