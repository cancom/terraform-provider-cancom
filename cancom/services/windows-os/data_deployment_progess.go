package windowsos

import (
	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataWindowsOSDeploymentProgress() *schema.Resource {
	return &schema.Resource{

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
	c := meta.(*client.Client)

	c.HostURL = c.ServiceURLs["managed-windows"]

	if d.Get("state").(string) == "Finished" {
		return nil
	}

	resp, err := (*client_windowsos.Client)(c).CreateWindowsDeploymentStatus(d.Get("deployment_id").(string))
	if err != nil {
		return err
	}

	d.SetId(resp.Id)
	d.Set("state", "Finished")

	return nil
}
