package windowsos

import (
	"context"
	"fmt"

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
	_, stateExists := d.GetOk("state")
	tflog.Debug(ctx, "Started Windows Deployment Progress")
	tflog.Debug(ctx, fmt.Sprintf("State is %b", stateExists))
	// if a status is already set, we can avoid calling the endpoint again.
	if (d.Get("state").(string)) == "Finished" {
		tflog.Debug(ctx, "Found already state finished")
		return nil
	}
	if (d.Get("state").(string)) == "Failed" {
		tflog.Debug(ctx, "Found already state failed")
		return nil
	}
	if (d.Get("state").(string)) == "Started" {
		tflog.Debug(ctx, "Found already state started")
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
		diags := diag.FromErr(err)

		//Change severity to waring to not break the deployment and store the values.
		for i := range diags {
			if diags[i].Severity == diag.Error {
				diags[i].Severity = diag.Warning
			}
		}
		return diags
	}
	d.SetId(resp.Id)
	d.Set("state", "Finished")
	return nil
}
