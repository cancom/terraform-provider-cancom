package windowsos

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWindowsOSDeploymentProgress() *schema.Resource {
	return &schema.Resource{

		Description: `Windows OS --- This data objects waits for CANCOM Manged Windows Cloud server deployment-system to complete. 
		
This can be used for dependency tracking purposes and returns in case of failure errors.  
The required` + " `deployment_id` " + `represents the id of` + " `cancom_windows_os_deployment` " + `resource.`,
		CreateContext: resourceWindowsOSDeploymentProgressCreate,
		ReadContext:   resourceWindowsOSDeploymentProgressRead,
		UpdateContext: resourceWindowsOSDeploymentProgressUpdate,
		DeleteContext: resourceWindowsOSDeploymentProgressDelete,
		Schema: map[string]*schema.Schema{
			"deployment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the deployment object",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the deployment",
			},
			"errorhandling": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "error",
				Description: `The errorhandling of deployment errors is configurable. If the output should be just warining and no interruption of the process is needed, please set the value to warning.

				!> For multi staged deployments the value error is highly recommended.`,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile("error|warning"), "errorhandling may only contain error or warning")),
			},
		},
	}
}

func resourceWindowsOSDeploymentProgressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorstatus := []int{6, 5}
	sucessstatus := []int{4}

	c, err := m.(*client.CcpClient).GetService("managed-windows")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	var deployment_id = d.Get("deployment_id").(string)
	d.SetId(deployment_id)
	d.Set("status", "started")
	for {
		resp, err := (*client_windowsos.Client)(c).GetWindowsDeploymentStatus(deployment_id)
		if err != nil {
			d.Set("status", "Failed")
			return diag.FromErr(err)
		}
		if slices.Contains(errorstatus, resp.Status) {
			err = fmt.Errorf("deployment for ID %s failed with statuscode %d", resp.Id, resp.Status)
			d.Set("status", "Failed")
			if d.Get("errorhandling").(string) == "warning" {
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Deployment failed - please review",
						Detail:   err.Error(),
					},
				}
			} else {
				return diag.FromErr(err)
			}

		}
		if slices.Contains(sucessstatus, resp.Status) {
			d.Set("status", "Succeded")
			return diags
		}
		time.Sleep(10 * time.Second) // sleep for 30 seconds to aviod active waiting
	}
}

func resourceWindowsOSDeploymentProgressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO: Discuss integration of evtl. and show a waring if the Deployment failed
	var diags diag.Diagnostics

	return diags
}

func resourceWindowsOSDeploymentProgressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Discuss integration of evtl. and show a waring if the Deployment failed
	var diags diag.Diagnostics

	return diags
}

func resourceWindowsOSDeploymentProgressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
