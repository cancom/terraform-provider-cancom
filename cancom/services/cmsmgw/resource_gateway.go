package cmsmgw

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_cmsmgw "github.com/cancom/terraform-provider-cancom/client/services/cmsmgw"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGateway() *schema.Resource {
	return &schema.Resource{
		Description:   "Managed Gateway --- Represents a gateway resource.",
		CreateContext: resourceGatewayCreate,
		ReadContext:   resourceGatewayRead,
		UpdateContext: resourceGatewayUpdate,
		DeleteContext: resourceGatewayDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bastion_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cancom_primary_gw_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cancom_secondary_gw_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mgw_size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_translation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"nat_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bastion_lite_linux": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bastion_lite_windows": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}

}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, "Creating Gateway")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	gateway := &client_cmsmgw.GatewayCreateRequest{
		Customer:           d.Get("customer").(string),
		NameTag:            d.Get("name_tag").(string),
		Tag:                d.Get("tag").(string),
		MgwSize:            d.Get("mgw_size").(string),
		NatTranslation:     d.Get("nat_translation").(bool),
		BastionLiteLinux:   d.Get("bastion_lite_linux").(bool),
		BastionLiteWindows: d.Get("bastion_lite_windows").(bool),
	}

	resp, err := (*client_cmsmgw.Client)(c).CreateGateway(gateway)

	if err != nil {
		return diag.FromErr(err)
	}
	mgwId := resp.ID

	//log.Println("[DEBUG] Something happened!%s, is %s", resp.ID, resp.Status)
	//resp.Status = "PENDING"
	i := 1
	for resp.State != "CREATE_COMPLETE" {
		time.Sleep(15 * time.Second)
		resp, err = (*client_cmsmgw.Client)(c).GetGateway(mgwId)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "Creating Gateway, Iteration: "+strconv.Itoa(i)+" Status: "+resp.State)

		if (resp.State != "CREATE_IN_PROGRESS") && (resp.State != "CREATE_COMPLETE") && (resp.State != "") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creationg mgw: %s", 401, resp.State)
			return diag.FromErr(err2)
		}

		i++
		if i >= 40 {
			//resp.Status = "UNKNOWN_ERROR"
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}

	d.SetId(resp.ID)

	resourceGatewayRead(ctx, d, m)

	return diags
}

// --------------Gateway Read----------------------------
func resourceGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//mgwId := d.Get("id").(string)
	mgwId := d.Id() // for import to work

	resp, err := (*client_cmsmgw.Client)(c).GetGateway(mgwId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	d.Set("customer", resp.Customer)
	d.Set("name_tag", resp.NameTag)
	d.Set("bastion_network", resp.BastionNetwork)
	d.Set("cancom_primary_gw_ip", resp.CancomPrimaryGwIp)
	d.Set("cancom_secondary_gw_ip", resp.CancomSecondaryGwIp)
	d.Set("state", resp.State)
	d.Set("nat_network", resp.NatNetwork)
	d.Set("nat_translation", resp.NatTranslation)
	d.Set("mgw_size", resp.MgwSize)
	d.Set("tag", resp.Tag)
	d.Set("bastion_lite_linux", resp.BastionLiteLinux)
	d.Set("bastion_lite_Window", resp.BastionLiteWindows)

	return diags
}

// ---------------update mgw----------------------
func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mgwId := d.Get("id").(string)

	gateway := &client_cmsmgw.GatewayUpdateRequest{
		ID:                 d.Get("id").(string),
		NameTag:            d.Get("name_tag").(string),
		Tag:                d.Get("tag").(string),
		MgwSize:            d.Get("mgw_size").(string),
		NatTranslation:     d.Get("nat_translation").(bool),
		BastionLiteLinux:   d.Get("bastion_lite_linux").(bool),
		BastionLiteWindows: d.Get("bastion_lite_windows").(bool),
	}

	resp, err := (*client_cmsmgw.Client)(c).UpdateGateway(mgwId, gateway)
	if err != nil {
		return diag.FromErr(err)
	}

	resp.State = "UPDATE_IN_PROGRESS"
	i := 1
	for (resp.State != "UPDATE_COMPLETE") && (resp.State != "CREATE_COMPLETE") {
		time.Sleep(15 * time.Second)
		resp, err = (*client_cmsmgw.Client)(c).GetGateway(mgwId)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "Updating Gateway, Iteration: "+strconv.Itoa(i)+" Status: "+resp.State)

		if (resp.State != "UPDATE_IN_PROGRESS") && (resp.State != "CREATE_COMPLETE") && (resp.State != "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS") && (resp.State != "UPDATE_COMPLETE") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creationg mgw: %s", 401, resp.State)
			return diag.FromErr(err2)
		}
		i++
		if i >= 40 {
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}
	resourceGatewayRead(ctx, d, m)
	return diags
}

// ----------------------delete mgw
func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mgwId := d.Get("id").(string)

	err = (*client_cmsmgw.Client)(c).DeleteGateway(mgwId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
