package akcauth

import (
	"context"
	"log"
	"terraform-provider-akcauth/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiResourceCreate,
		ReadContext:   resourceApiResourceRead,
		UpdateContext: resourceApiResourceUpdate,
		DeleteContext: resourceApiResourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceApiResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	displayName := d.Get("display_name").(string)
	scopesRaw := d.Get("scopes").([]interface{})
	scopes := make([]string, len(scopesRaw))
	for i, raw := range scopesRaw {
		scopes[i] = raw.(string)
	}

	model := client.ApiResourceCreate{
		DisplayName: displayName,
		Scopes:      scopes,
	}

	err := c.CreateApiResource(name, &model)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	resourceApiResourceRead(ctx, d, m)

	return diags
}

func resourceApiResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	apiResourceName := d.Id()

	authApiResource, err := c.GetApiResource(apiResourceName)
	if err != nil {
		log.Printf("[WARN] Api resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("name", authApiResource.Name)
	d.Set("display_name", authApiResource.DisplayName)
	d.Set("scopes", authApiResource.Scopes)

	return diags
}

func resourceApiResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	name := d.Id()

	apiResource, err := c.GetApiResource(name)
	if err != nil {
		_, ok := err.(*client.ClientError)
		if ok {
			return diag.FromErr(err)
		} else {
			return diag.FromErr(err)
		}
	}

	updateModel := apiResource.ToUpdateModel()

	if d.HasChange("display_name") {
		updateModel.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("scopes") {
		scopesRaw := d.Get("scopes").([]interface{})
		scopes := make([]string, len(scopesRaw))
		for i, raw := range scopesRaw {
			scopes[i] = raw.(string)
		}
		updateModel.Scopes = scopes
	}

	err = c.UpdateApiResource(name, &updateModel)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApiResourceRead(ctx, d, m)
}

func resourceApiResourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	name := d.Id()

	err := c.DeleteApiResource(name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
