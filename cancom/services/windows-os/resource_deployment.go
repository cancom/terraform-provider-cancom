package windowsos

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_windowsos "github.com/cancom/terraform-provider-cancom/client/services/windows-os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWindowsOSDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWindowsOSCreate,
		ReadContext:   resourceWindowsOSRead,
		UpdateContext: resourceWindowsOSUpdate,
		DeleteContext: resourceWindowsOSDelete,
		Schema: map[string]*schema.Schema{
			"customer_environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the customer environment.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role of the Machine - needs to match predefined roles.",
			},
			"maintenance_window_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of maintenanceWindow IDs",
			},
			"services": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Services deployed and delivered to maschine",
			},
			"customer_number": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				Description: "Customer number for Windows Management",
			},
		},
	}
}

func resourceWindowsOSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["managed-windows"]

	var diags diag.Diagnostics

	tempservices := d.Get("services").([]interface{})
	servicesarray := []string{}
	for _, tempservices := range tempservices {
		servicesarray = append(servicesarray, tempservices.(string))
	}

	tempmaintenance_window_id := d.Get("maintenance_window_id").([]interface{})
	maintenance_window_id_array := []string{}
	for _, tempmaintenance_window_id := range tempmaintenance_window_id {
		maintenance_window_id_array = append(maintenance_window_id_array, tempmaintenance_window_id.(string))
	}
	computerObject := client_windowsos.WindowsOS_Computer{
		Services:            servicesarray,
		MaintenanceWindowId: maintenance_window_id_array,
		Role:                d.Get("role").(string),
	}

	windowsOSRequest := client_windowsos.WindowsOS_Deplyoment{

		Computer:               computerObject,
		CoustomerEnvironmentID: d.Get("customer_environment_id").(string),
	}

	resp, err := (*client_windowsos.Client)(c).CreateWindowsDeployment(&windowsOSRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("customer_number", resp.Computer.CustomerID)
	d.SetId(resp.Id)

	return diags

}

func resourceWindowsOSRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["managed-windows"]

	var diags diag.Diagnostics

	resp, err := (*client_windowsos.Client)(c).GetWindowsDeployment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("customer_environment_id", resp.CoustomerEnvironmentID)
	d.Set("role", resp.Computer.Role)
	d.Set("maintenance_window_id", resp.Computer.MaintenanceWindowId)
	d.Set("services", resp.Computer.Services)

	return diags
}

func resourceWindowsOSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["managed-windows"]

	var diags diag.Diagnostics


	tempservices := d.Get("services").([]interface{})
	servicesarray := []string{}
	for _, tempservices := range tempservices {
		servicesarray = append(servicesarray, tempservices.(string))
	}

	tempmaintenance_window_id := d.Get("maintenance_window_id").([]interface{})
	maintenance_window_id_array := []string{}
	for _, tempmaintenance_window_id := range tempmaintenance_window_id {
		maintenance_window_id_array = append(maintenance_window_id_array, tempmaintenance_window_id.(string))
	}
	computerObject := client_windowsos.WindowsOS_Computer{
		Services:            servicesarray,
		MaintenanceWindowId: maintenance_window_id_array,

		Role:                d.Get("role").(string),
	}

	windowsOSRequest := client_windowsos.WindowsOS_Deplyoment{
		Computer:               computerObject,
		CoustomerEnvironmentID: d.Get("customer_environment_id").(string),
	}

	resp, err := (*client_windowsos.Client)(c).UpdateWindowsOsDeployment(d.Id(), &windowsOSRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Id)

	return diags
}

func resourceWindowsOSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["managed-windows"]

	var diags diag.Diagnostics

	err := (*client_windowsos.Client)(c).DeleteWindowsDeployment(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
