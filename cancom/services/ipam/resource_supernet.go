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
		Description:   "IP Management --- IPAM Supernets are larger aggregates that are further subnetted by the service.",
		CreateContext: resourceSupernetCreate,
		ReadContext:   resourceSupernetRead,
		UpdateContext: resourceSupernetUpdate,
		DeleteContext: resourceSupernetDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
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
			"supernet_cidr": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		StateUpgraders: []schema.StateUpgrader{
			supernetUpgradeV0(),
		},
	}
}

func resourceSupernetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := (*client_ipam.Client)(c).GetSupernet(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name_tag", resp.NameTag)
	d.Set("instance_id", resp.InstanceId)
	d.Set("supernet_cidr", resp.SupernetCidr)
	d.Set("description", resp.Description)
	d.Set("created_at", resp.CreatedAt)
	d.Set("updated_at", resp.UpdatedAt)
	d.Set("id", resp.ID)

	return diags
}

func resourceSupernetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	supernet := &client_ipam.SupernetCreateRequest{
		InstanceId:   d.Get("instance_id").(string),
		NameTag:      d.Get("name_tag").(string),
		Description:  d.Get("description").(string),
		SupernetCidr: d.Get("supernet_cidr").(string),
	}

	resp, err := (*client_ipam.Client)(c).CreateSupernet(supernet)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(id)

	resourceSupernetRead(ctx, d, m)

	return diags

}

func resourceSupernetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	supernet := &client_ipam.SupernetUpdateRequest{
		NameTag:      d.Get("name_tag").(string),
		SupernetCidr: d.Get("supernet_cidr").(string),
		InstanceId:   d.Get("instance_id").(string),
		Description:  d.Get("description").(string),
	}

	_, err = (*client_ipam.Client)(c).UpdateSupernet(id, supernet)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceSupernetRead(ctx, d, m)

	return diags
}

func resourceSupernetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	err = (*client_ipam.Client)(c).DeleteSupernet(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
