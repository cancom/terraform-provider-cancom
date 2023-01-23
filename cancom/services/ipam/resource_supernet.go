package ipam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_ipam "github.com/cancom/terraform-provider-cancom/client/services/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSupernet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSupernetCreate,
		// ReadContext:   resourceSupernetRead,
		// UpdateContext: resourceSupernetUpdate,
		// DeleteContext: resourceSupernetDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the supernet",
				Required:    true,
			},
		},
	}
}

func resourceSupernetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ip-management"]

	var diags diag.Diagnostics

	supernet := &client_ipam.SupernetCreateRequest{
		Name: d.Get("name").(string),
	}

	resp, err := (*client_ipam.Client)(c).CreateSupernet(supernet)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.Set("name", resp.Name)

	d.SetId(id)

	// resourceGatewayRead(ctx, d, m)

	return diags

}
