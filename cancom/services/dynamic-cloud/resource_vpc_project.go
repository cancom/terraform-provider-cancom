package dynamiccloud

import (
	"context"
	"regexp"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dynamiccloud "github.com/cancom/terraform-provider-cancom/client/services/dynamic-cloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVpcProject() *schema.Resource {
	nameRegex := regexp.MustCompile("[a-zA-Zß0-9-_]+")
	commentRegex := regexp.MustCompile("[a-zA-Zß0-9-_.,;:?!#+ ]+")
	return &schema.Resource{
		CreateContext: resourceVpcProjectCreate,
		ReadContext:   resourceVpcProjectRead,
		DeleteContext: resourceVpcProjectDelete,
		Schema: map[string]*schema.Schema{
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the user who created the VPC Project.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the date when the VPC Project was created.",
			},
			"limits": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The resource limits currently configured for this VPC Project.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user defined name used to construct the OpenStack project name with the schema '\\<tenant\\>-\\<name\\>'.</br>By changing this value the old project will be delete and a new project with the new name will be created.</br> ~> **WARNING:** Changing this value will delete all resources in the VPC Project.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 41),
					validation.StringMatch(nameRegex, "Name may only contain (a-zA-Zß0-9-_)."),
				)),
			},
			"openstack_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uuid of the OpenStack Project.",
			},
			"project_comment": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          "",
				Description:      "A comment to describe what this VPC Project is used for.</br>By changing this value the old project will be delete and a new project will be created.</br> ~> **WARNING:** Changing this value will delete all resources in the VPC Project.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(commentRegex, "project_comment may only contain (a-zA-Zß0-9-_.,;:?!#+ ).")),
			},
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the tenant this VPC Project belongs to.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceVpcProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating VPC Project")
	vpcProjectRequest := client_dynamiccloud.VpcProjectCreateRequest{
		Metadata: client_dynamiccloud.VpcProjectCreateMetadata{
			Name: d.Get("name").(string),
		},
		Spec: client_dynamiccloud.VpcProjectCreateSpec{
			ProjectComment: d.Get("project_comment").(string),
		},
	}

	resp, err := (*client_dynamiccloud.Client)(c).CreateVpcProject(&vpcProjectRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	vpcProjectShortid := resp.Metadata.Shortid
	d.SetId(vpcProjectShortid)

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {
		return WaitProjectReady(ctx, c, vpcProjectShortid)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcProjectRead(ctx, d, meta)
}

func resourceVpcProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Reading VPC Project")
	var diags diag.Diagnostics
	vpcProjectShortid := d.Id()

	// GetVpcProject returns nil if the VPC Project is NotFound
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		d.SetId("")
		return nil
	}

	d.Set("created_by", resp.Spec.CreatedBy)
	d.Set("creation_date", resp.Metadata.CreationDate)
	err = d.Set("limits", resp.Spec.Limits)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", resp.Metadata.Name)
	d.Set("openstack_uuid", resp.Spec.OpenstackUuid)
	d.Set("project_comment", resp.Spec.ProjectComment)
	d.Set("tenant", resp.Metadata.Tenant)

	return diags
}

func resourceVpcProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Deleting VPC Project")
	var diags diag.Diagnostics
	vpcProjectShortid := d.Id()

	err = (*client_dynamiccloud.Client)(c).DeleteVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete)-time.Minute, func() *resource.RetryError {
		return WaitProjectDeleted(ctx, c, vpcProjectShortid)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
