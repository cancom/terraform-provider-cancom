package object_storage

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/cancom/terraform-provider-cancom/client"
	client_object_storage "github.com/cancom/terraform-provider-cancom/client/services/object-storage"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description: `Object Storage --- Create a new user for the S3 API and assign permissions to the user.

You can set IAM policy documents that adheres to the AWS Identity and Access Management (IAM) policy structure. 
The document defines permissions using a JSON-based format, following the official [AWS specifications](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies.html) for IAM policies. 

Example: ` + "\n\n```json" + `
{
	"Statement": [
			{
				"Effect": "Allow",
				"Action": "*",
				"Resource": "*"
			}
	]
}
` + "```" + `

We recommend using the ` + "`jsonencode`" + ` TerraForm function to convert your policy document into a JSON string.`,
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		UpdateContext: resourceUserUpdate,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the user. The name must be unique inside your account.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description what the user is used for",
			},
			"permissions": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The IAM permission document.",
				ValidateDiagFunc: permissionValidator,
			},
		},
	}
}

func permissionValidator(data interface{}, path cty.Path) diag.Diagnostics {
	var policy map[string]interface{}
	err := json.Unmarshal([]byte(data.(string)), &policy)
	if err != nil {
		return diag.FromErr(err)
	}

	statement, ok := policy["Statement"]
	if !ok {
		return diag.Errorf("the policy must contain a statement")
	}

	if reflect.TypeOf(statement).Kind() != reflect.Slice {
		return diag.Errorf("the statement in the permission policy must be an array")
	}

	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	resp, err := (*client_object_storage.Client)(c).GetUser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := json.Marshal(resp.Permissions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("username", resp.Username)
	d.Set("description", resp.Description)
	d.Set("permissions", policy)

	return diags
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	var policy map[string]interface{}
	err = json.Unmarshal([]byte(d.Get("permissions").(string)), &policy)
	if err != nil {
		return diag.FromErr(err)
	}

	userCreateRequest := client_object_storage.UserCreateRequest{
		Username:    d.Get("username").(string),
		Description: d.Get("description").(string),
		Permissions: policy,
	}

	resp, err := (*client_object_storage.Client)(c).CreateUser(&userCreateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Id)
	d.Set("username", resp.Username)
	d.Set("description", resp.Description)
	d.Set("permissions", resp.Permissions)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	err = (*client_object_storage.Client)(c).DeleteUser(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	id := d.Id()

	var policy map[string]interface{}
	err = json.Unmarshal([]byte(d.Get("permissions").(string)), &policy)
	if err != nil {
		return diag.FromErr(err)
	}

	user := &client_object_storage.UserUpdateRequest{
		Description: d.Get("description").(string),
		Permissions: policy,
	}

	resp, err := (*client_object_storage.Client)(c).UpdateUser(id, user)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	d.Set("username", resp.Username)
	d.Set("description", resp.Description)
	d.Set("permissions", resp.Permissions)

	return diags
}
