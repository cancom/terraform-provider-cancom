package cmsmgw

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_cmsmgw "github.com/cancom/terraform-provider-cancom/client/services/cmsmgw"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatewayCreate,
		// ReadContext:   resourceGatewayRead,
		// UpdateContext: resourceGatewayUpdate,
		// DeleteContext: resourceGatewayDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the gateway",
				Required:    true,
			},
		},
	}
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["cmsmgw"]

	var diags diag.Diagnostics

	record := &client_cmsmgw.GatewayCreateRequest{
		Name: d.Get("name").(string),
	}

	resp, err := (*client_cmsmgw.Client)(c).CreateGateway(record)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.Set("name", resp.Name)

	d.SetId(id)

	// resourceGatewayRead(ctx, d, m)

	return diags

}
