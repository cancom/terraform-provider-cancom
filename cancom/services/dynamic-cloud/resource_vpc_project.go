package dynamiccloud

import (
	"context"
	"fmt"
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
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the date when the VPC Project was created.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user customizable name which will be used to construct the OpenStack projects name with the schema '<tenant>-<name>'.\nBy changing this value the old project will be delete and a new project with the new name will be created.\nWARNING: Recreation will delete all resources in the VPC Project.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 41),
					validation.StringMatch(nameRegex, "Name may only contain (a-zA-Zß0-9-_)."),
				)),
			},
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the tenant this VPC Project belongs to.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the user who created the VPC Project.",
			},
			"project_comment": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true, //TODO change this to false once backend supports changing the comment
				Default:          "",
				Description:      "A comment to describe what this VPC Project is used for.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(commentRegex, "project_comment may only contain (a-zA-Zß0-9-_.,;:?!#+ ).")),
			},
			"openstack_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uuid of the OpenStack Project.",
			},
			"limits": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The resource limits currently configured for this VPC Project.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The phase the VPC Project is in, possible values are `Creating`, `Ready`, `Updating`, `Deleting` and `Error`.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceVpcProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("dynamic-cloud")
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
		// GetVpcProject returns nil if the VPC Project is NotFound
		resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing VPC Project during creation: %s", err))
		}
		if resp == nil {
			return resource.NonRetryableError(fmt.Errorf("error describing VPC Project during creation. VPC Project NotFound"))
		}

		if resp.Status.Phase != "Ready" {
			tflog.Info(ctx, "Waiting for VPC Project status phase to become Ready")
			return resource.RetryableError(fmt.Errorf("expected VPC Project status phase to be Ready but was in phase %s", resp.Status.Phase))
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcProjectRead(ctx, d, m)
}

func resourceVpcProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	vpcProjectShortid := d.Id()

	// GetVpcProject returns nil if the VPC Project is NotFound
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.FromErr(fmt.Errorf("error describing VPC Project. VPC Project NotFound"))
	}

	d.Set("creation_date", resp.Metadata.CreationDate)
	d.Set("name", resp.Metadata.Name)
	d.Set("tenant", resp.Metadata.Tenant)
	d.Set("created_by", resp.Spec.CreatedBy)
	d.Set("project_comment", resp.Spec.ProjectComment)
	d.Set("openstack_uuid", resp.Spec.OpenstackUuid)
	d.Set("limits", resp.Spec.Limits)
	d.Set("phase", resp.Status.Phase)

	return diags
}

func resourceVpcProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("dynamic-cloud")
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
		// GetVpcProject returns nil if the VPC Project is NotFound
		resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing VPC Project during deletion: %s", err))
		}

		// resp not nil, so VPC Project deletion not yet finished / VPC Project not yet NotFound
		if resp != nil {
			tflog.Info(ctx, "Waiting for VPC Project to be Deleted")
			return resource.RetryableError(fmt.Errorf("expected VPC Project to be Deleted but was in state %s", resp.Status.Phase))
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
