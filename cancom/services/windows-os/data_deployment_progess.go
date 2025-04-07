package windowsos

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWindowsOSDeploymentProgress() *schema.Resource {
	return &schema.Resource{
		Description: `Windows OS --- This data objects waits for CANCOM Manged Windows Cloud server deployment-system to complete. 
		
This can be used for dependency tracking purposes and returns in case of failure errors.  
The required` + " `deployment_id` " + `represents the id of` + " `cancom_windows_os_deployment` " + `resource.`,
		ReadContext: WindowsOSDeploymentProgressRead,
		Schema: map[string]*schema.Schema{
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the deployment object.",
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func WindowsOSDeploymentProgressRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("managed-windows")
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Started Windows Deployment Progress")

	// if a status is already set, we can avoid calling the endpoint again.
	if d.Get("state").(string) == "Finished" {
		return nil
	} else if d.Get("state").(string) == "Failed" {
		return nil
	} else if d.Get("state").(string) == "Started" {
		return nil
	}
	tflog.Debug(ctx, "Finished Pre-Check")

	d.SetId(d.Get("deployment_id").(string))
	d.Set("state", "Started")

	tflog.Debug(ctx, "Updated State")
	resp, err := (*client_windowsos.Client)(c).CreateWindowsDeploymentStatus(d.Get("deployment_id").(string))
	if err != nil {
		d.SetId(d.Get("deployment_id").(string))
		d.Set("state", "Failed")
		return diag.FromErr(err)
	}
	d.SetId(resp.Id)
	d.Set("state", "Finished")

	return nil
}
