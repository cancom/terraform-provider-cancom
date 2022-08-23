package iam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
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

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	userCreateRequest := client_iam.UserCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	resp, err := (*client_iam.Client)(c).CreateUser(&userCreateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceUserRead(ctx, d, m)

	d.SetId(resp.Principal)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	user := client_iam.User{
		Principal: d.Id(),
	}

	resp, err := (*client_iam.Client)(c).GetUser(user.Principal)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", resp.Description)
	d.Set("group", resp.Group)
	d.Set("principal", resp.Principal)
	d.Set("tenant", resp.Tenant)
	d.Set("created_at", resp.CreatedAt)
	d.Set("created_by", resp.CreatedBy)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	user := client_iam.User{
		Principal: d.Id(),
	}

	userUpdateRequest := client_iam.UserUpdateRequest{
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	err := (*client_iam.Client)(c).UpdateUser(user.Principal, &userUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	user := client_iam.User{
		Principal: d.Id(),
	}

	err := (*client_iam.Client)(c).DeleteUser(user.Principal)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
