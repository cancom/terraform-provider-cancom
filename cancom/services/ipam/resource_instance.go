package ipam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_ipam "github.com/cancom/terraform-provider-cancom/client/services/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_tag": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"managed_by": {
				Type:     schema.TypeString,
				Computed: false,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: false,
				Optional: true,
			},
			"release_wait_time": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	resp, err := (*client_ipam.Client)(c).GetInstance(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name_tag", resp.NameTag)
	d.Set("managed_by", resp.ManagedBy)
	d.Set("description", resp.Description)
	d.Set("release_wait_time", resp.ReleaseWaitTime)
	d.Set("created_at", resp.CreatedAt)
	d.Set("updated_at", resp.UpdatedAt)
	d.Set("id", resp.ID)

	return diags
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	instance := &client_ipam.InstanceCreateRequest{
		NameTag:         d.Get("name_tag").(string),
		ManagedBy:       d.Get("managed_by").(string),
		ReleaseWaitTime: d.Get("release_wait_time").(string),
		Description:     d.Get("description").(string),
	}

	resp, err := (*client_ipam.Client)(c).CreateInstance(instance)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(id)

	resourceInstanceRead(ctx, d, m)

	return diags

}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	instance := &client_ipam.InstanceUpdateRequest{
		NameTag:         d.Get("name_tag").(string),
		ManagedBy:       d.Get("managed_by").(string),
		ReleaseWaitTime: d.Get("release_wait_time").(string),
		Description:     d.Get("description").(string),
	}

	_, err = (*client_ipam.Client)(c).UpdateInstance(id, instance)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceInstanceRead(ctx, d, m)

	return diags
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("ip-management")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	err = (*client_ipam.Client)(c).DeleteInstance(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
