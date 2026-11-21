package iam

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "IAM --- User with long-term credentials for automation purpose.",
		CreateContext: resourceServiceUserCreate,
		ReadContext:   resourceServiceUserRead,
		UpdateContext: resourceServiceUserUpdate,
		DeleteContext: resourceServiceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the service user. The name must be unique in your account.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for what the service user is used for.",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group that the service user belongs to.",
			},
			"principal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the service user.",
			},
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant that the service user belongs to.",
			},
			"session": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Session token of the service user. Only returned on creation.",
			},
			"session_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deprecated. Must not be used.",
				Deprecated:  "`session_hash` is deprecated and might be removed in a future version.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The initial creation date of the service user.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The initial creator of the service user.",
			},
		},
	}
}

func resourceServiceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	userCreateRequest := client_iam.ServiceUserCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	serviceUser, err := (*client_iam.Client)(c).CreateServiceUser(&userCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating user: %s", err)
	}

	d.SetId(serviceUser.Principal)
	if err := d.Set("session", serviceUser.Session); err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceUserRead(ctx, d, m)
}

func resourceServiceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	serviceUser, err := (*client_iam.Client)(c).GetServiceUser(d.Id())
	if err != nil {
		return diag.Errorf("Error reading user: %s", err)
	}

	d.Set("description", serviceUser.Description)
	d.Set("group", serviceUser.Group)
	d.Set("principal", serviceUser.Principal)
	d.Set("tenant", serviceUser.Tenant)
	d.Set("session_hash", serviceUser.SessionHash)
	d.Set("created_at", serviceUser.CreatedAt)
	d.Set("created_by", serviceUser.CreatedBy)

	return diags
}

func resourceServiceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	serviceUserUpdateRequest := client_iam.ServiceUserUpdateRequest{
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	err = (*client_iam.Client)(c).UpdateServiceUser(d.Id(), &serviceUserUpdateRequest)
	if err != nil {
		return diag.Errorf("Error updating user: %s", err)
	}

	return resourceServiceUserRead(ctx, d, m)
}

func resourceServiceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	err = (*client_iam.Client)(c).DeleteServiceUser(d.Id())
	if err != nil {
		return diag.Errorf("Error deleting user: %s", err)
	}

	d.SetId("")

	return diags
}
