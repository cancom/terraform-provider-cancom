package ipam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_ipam "github.com/cancom/terraform-provider-cancom/client/services/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "IP Management --- IPAM networks are assigned by the service from supernets.",
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supernet_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
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
			"request": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"host_assign": {
				Type:     schema.TypeBool,
				Computed: false,
				Optional: true,
				Default:  true,
			},
			"prefix_str": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := (*client_ipam.Client)(c).GetNetwork(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name_tag", resp.NameTag)
	d.Set("supernet_id", resp.SupernetId)
	d.Set("description", resp.Description)
	d.Set("prefix_str", resp.PrefixStr)
	d.Set("host_assign", resp.HostAssign)
	d.Set("created_at", resp.CreatedAt)
	d.Set("updated_at", resp.UpdatedAt)
	d.Set("id", resp.ID)

	return diags
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	network := &client_ipam.NetworkCreateRequest{
		SupernetId:  d.Get("supernet_id").(string),
		Request:     d.Get("request").(string),
		NameTag:     d.Get("name_tag").(string),
		Description: d.Get("description").(string),
	}

	resp, err := (*client_ipam.Client)(c).CreateNetwork(network)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(id)

	// if hostAssign is true we must enable ist in an update to th network we just assigned
	if !d.Get("host_assign").(bool) {
		update := &client_ipam.NetworkUpdateRequest{
			HostAssign: false,
			//Source:    "CANCOM-TF",
		}
		_, update_err := (*client_ipam.Client)(c).UpdateNetwork(id, update)
		if update_err != nil {
			return diag.FromErr(update_err)
		}
	}

	resourceNetworkRead(ctx, d, m)

	return diags

}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	network := &client_ipam.NetworkUpdateRequest{
		NameTag:     d.Get("name_tag").(string),
		Description: d.Get("description").(string),
		HostAssign:  d.Get("host_assign").(bool),
	}

	_, err = (*client_ipam.Client)(c).UpdateNetwork(id, network)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceNetworkRead(ctx, d, m)

	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	// first try to disable host assignment if it is enabled.
	network := &client_ipam.NetworkUpdateRequest{
		HostAssign: false,
		//Source:    "CANCOM-TF",
	}
	_, err = (*client_ipam.Client)(c).UpdateNetwork(id, network)
	if err != nil {
		return diag.FromErr(err)
	}

	err = (*client_ipam.Client)(c).DeleteNetwork(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
