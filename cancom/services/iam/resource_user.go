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
		Description:   "IAM --- IAM user authenticated against the central CANCOM SSO.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the user. The name must be unique in your account.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for what the user is used for.",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group that the user belongs to.",
			},
			"principal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the principal.",
			},
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant that the user belongs to.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The initial creation date of the user.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The initial creator of the user.",
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

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
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

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
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	user := client_iam.User{
		Principal: d.Id(),
	}

	userUpdateRequest := client_iam.UserUpdateRequest{
		Description: d.Get("description").(string),
		Group:       d.Get("group").(string),
	}

	err = (*client_iam.Client)(c).UpdateUser(user.Principal, &userUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	user := client_iam.User{
		Principal: d.Id(),
	}

	err = (*client_iam.Client)(c).DeleteUser(user.Principal)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
