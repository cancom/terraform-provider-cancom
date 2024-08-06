package windowsos

import (
	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWindowsOSDeploymentProgress() *schema.Resource {
	return &schema.Resource{
		Description: `Windows OS --- This data objects waits for CANCOM Manged Windows Cloud server deployment-system to complete. 
		
This can be used for dependency tracking purposes and returns in case of failure errors.  
The required` + " `deployment_id` " + `represents the id of` + " `cancom_windows_os_deployment` " + `resource.`,
		Read: WindowsOSDeploymentProgressRead,
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

func WindowsOSDeploymentProgressRead(d *schema.ResourceData, meta interface{}) error {
	c, err := meta.(*client.CcpClient).GetService("managed-windows")
	if err != nil {
		return err
	}
	// if a status is already set, we can avoid calling the endpoint again.
	if d.Get("state").(string) == "Finished" {
		return nil
	} else if d.Get("state").(string) == "Failed" {
		return nil
	}

	resp, err := (*client_windowsos.Client)(c).CreateWindowsDeploymentStatus(d.Get("deployment_id").(string))
	if err != nil {
		d.SetId(d.Get("deployment_id").(string))
		d.Set("state", "Failed")
		return err
	}

	d.SetId(resp.Id)
	d.Set("state", "Finished")

	return nil
}
