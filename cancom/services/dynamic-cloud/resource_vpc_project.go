package dynamiccloud

import (
	"context"
	"fmt"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dynamiccloud "github.com/cancom/terraform-provider-cancom/client/services/dynamic-cloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpcProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcProjectCreate,
		ReadContext:   resourceVpcProjectRead,
		UpdateContext: resourceVpcProjectUpdate,
		DeleteContext: resourceVpcProjectDelete,
		Schema: map[string]*schema.Schema{
			"metadata_creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creating date of the VPC Project",
			},
			"metadata_deletion_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Deletion date of the VPC Project",
			},
			"metadata_name": {
				Type:        schema.TypeString,
				Optional:    false,
				ForceNew:    true,
				Description: "The display name which will be used as suffix for the OpenStack Project name",
			},
			"metadata_shortid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of the VPC Project",
			},
			"metadata_tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The tenant this VPC Project belongs to",
			},
			"spec_created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of the user who created the VPC Project.",
			},
			"spec_project_comment": {
				Type:        schema.TypeString,
				Required:    false,
				ForceNew:    false,
				Description: "Comment may be used to describe what this VPC Project is used for.",
			},
			"spec_openstack_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uuid of the OpenStack Project.",
			},
			"spec_limits": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The resource limits currently configured for this VPC Project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec_project_users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The crns of the users to get member permissions in the VPC Project.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"spec_managed_by_service": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The tenant this VPC Project belongs to",
			},
			"status_phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The phase the VPC Project is in, possible values are `Creating`, `Ready`, `Updating`, `Deleting` and `Error`.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
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
		Name:           d.Get("metadata_name").(string),
		ProjectComment: d.Get("spec_project_comment").(string),
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

	d.Set("metadata_creation_date", resp.Metadata.CreationDate)
	d.Set("metadata_deletion_date", resp.Metadata.DeletionDate)
	d.Set("metadata_name", resp.Metadata.Name)
	d.Set("metadata_shortid", resp.Metadata.Shortid)
	d.Set("metadata_tenant", resp.Metadata.Tenant)
	d.Set("spec_created_by", resp.Spec.CreatedBy)
	d.Set("spec_project_comment", resp.Spec.ProjectComment)
	d.Set("spec_openstack_uuid", resp.Spec.OpenstackUuid)
	d.Set("spec_limits", resp.Spec.Limits)
	d.Set("spec_project_users", resp.Spec.ProjectUsers)
	d.Set("spec_managed_by_service", resp.Spec.ManagedByService)
	d.Set("status_phase", resp.Status.Phase)

	return diags
}

func resourceVpcProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
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
