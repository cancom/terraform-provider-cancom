package ipam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_ipam "github.com/cancom/terraform-provider-cancom/client/services/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_crn": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"qualifier": {
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"name_tag": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	//c.HostURL = c.ServiceURLs["ip-management"]
	c.HostURL = "https://ip-management.portal.cancom.io"
	var diags diag.Diagnostics

	id := d.Id()

	resp, err := (*client_ipam.Client)(c).GetHost(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name_tag", resp.NameTag)
	d.Set("address", resp.Address)
	d.Set("description", resp.Description)
	d.Set("id", resp.ID)

	return diags
}

func resourceHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ip-management"]

	var diags diag.Diagnostics

	host := &client_ipam.HostCreateRequest{
		NetworkCrn:  d.Get("network_crn").(string),
		NameTag:     d.Get("name_tag").(string),
		Operation:   "assign_address",
		Description: d.Get("description").(string),
		Qualifier:   d.Get("qualifier").(string),
		Source:      "CANCOM-TF",
	}

	resp, err := (*client_ipam.Client)(c).CreateHost(host)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.Set("address", resp.Address)
	d.SetId(id)

	// resourceGatewayRead(ctx, d, m)

	return diags

}

func resourceHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ip-management"]

	var diags diag.Diagnostics

	id := d.Id()

	host := &client_ipam.HostUpdateRequest{
		NetworkCrn:  d.Get("network_crn").(string),
		NameTag:     d.Get("name_tag").(string),
		Description: d.Get("description").(string),
		Source:      "CANCOM-TF",
	}

	_, err := (*client_ipam.Client)(c).UpdateHost(id, host)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceHostRead(ctx, d, m)

	return diags
}

func resourceHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ip-management"]

	var diags diag.Diagnostics

	id := d.Id()

	err := (*client_ipam.Client)(c).DeleteHost(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
