package iam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"principal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	roleCreateRequest := &client_iam.RoleCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	role, err := (*client_iam.Client)(c).CreateRole(roleCreateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceRoleRead(ctx, d, m)

	d.SetId(role.Principal)

	return diags

}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	role, err := (*client_iam.Client)(c).GetRole(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", role.Description)
	d.Set("group", role.Group)
	d.Set("principal", role.Principal)
	d.Set("tenant", role.Tenant)
	d.Set("created_at", role.CreatedAt)
	d.Set("created_by", role.CreatedBy)

	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	roleUpdateRequest := &client_iam.RoleUpdateRequest{
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	err := (*client_iam.Client)(c).UpdateRole(d.Id(), roleUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceRoleRead(ctx, d, m)

	return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	err := (*client_iam.Client)(c).DeleteRole(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
