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

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
		Schema: map[string]*schema.Schema{
			// id is the id of the mgw, must be provied
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mgw_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deployment_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// this is the gateway ip on the customer side, cancom side ist given in the mgw resource
			"customer_primary_gw_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			// this is the gateway ip on the customer side, cancom side ist given in the mgw resource
			"customer_secondary_gw_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_profile": {
				Type:     schema.TypeString,
				Required: true,
			},
			// if this is not provided the backend will generate a secret
			"ipsec_psk_a": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			// if this is not provided the backend will generate a secret
			"ipsec_psk_b": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"cancom_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"customer_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("cmsmgw")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	temp1 := d.Get("cancom_networks").([]interface{})
	var cancomNetworks []string
	for _, temp1 := range temp1 {
		cancomNetworks = append(cancomNetworks, temp1.(string))
	}

	temp2 := d.Get("customer_networks").([]interface{})
	var customerNetworks []string
	for _, temp2 := range temp2 {
		customerNetworks = append(customerNetworks, temp2.(string))
	}

	connection := &client_cmsmgw.ConnectionCreateRequest{
		MgwId:                 d.Get("mgw_id").(string),
		NameTag:               d.Get("name_tag").(string),
		CustomerPrimaryGwIp:   d.Get("customer_primary_gw_ip").(string),
		CustomerSecondaryGwIp: d.Get("customer_secondary_gw_ip").(string),
		ConnectionProfile:     d.Get("connection_profile").(string),
		IpsecPskA:             d.Get("ipsec_psk_a").(string),
		IpsecPskB:             d.Get("ipsec_psk_b").(string),
		CancomNetworks:        cancomNetworks,
		CustomerNetworks:      customerNetworks,
	}

	resp, err := (*client_cmsmgw.Client)(c).CreateConnection(connection)
	if err != nil {
		return diag.FromErr(err)
	}
	connectionId := resp.ID
	// this is to give dynamo time to get consistant
	//time.Sleep(60 * time.Second)
	resp.Status = "PENDING_DEPLOYMENT"
	err_count := 1
	i := 1
	for resp.Status != "DEPLOYED" {
		time.Sleep(15 * time.Second)
		resp_new, err := (*client_cmsmgw.Client)(c).GetConnection(connectionId)
		if err != nil {
			err_count++
			if err_count > 6 {
				return diag.FromErr(err)
			}
			tflog.Info(ctx, "Creating Connection, skipping error on iteration:  "+strconv.Itoa(i)+" Status: "+resp.Status)
			time.Sleep(30 * time.Second)
		} else {
			resp = resp_new
		}
		tflog.Info(ctx, "Creating Connection, Iteration: "+strconv.Itoa(i)+" Status: "+resp.Status)

		if (resp.Status != "PENDING_DEPLOYMENT") && (resp.Status != "DEPLOYED") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creationg mgw: %s", 404, resp.Status)
			return diag.FromErr(err2)
		}

		/* for debugging
		if resp.Status == "CREATE_COMPLETE" {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Debug, Create complete found: %s", 404, resp.Status)
			return diag.FromErr(err2)
		}*/

		i++
		if i >= 40 {
			resp.Status = "UNKNOWN_ERROR"
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}

	d.SetId(resp.ID)

	resourceConnectionRead(ctx, d, m)

	return diags
}

// --------------------------Read Connection------------------------------
func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("cmsmgw")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	connectionId := d.Get("id").(string)

	//order, err := c.GetOrder(orderId)
	resp, err := (*client_cmsmgw.Client)(c).GetConnection(connectionId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	d.Set("mgw_id", resp.MgwId)
	d.Set("name_tag", resp.NameTag)
	d.Set("customer_primary_gw_ip", resp.CustomerPrimaryGwIp)
	d.Set("customer_secondary_gw_ip", resp.CustomerSecondaryGwIp)
	d.Set("deployment_state", resp.Status)
	d.Set("connection_profile", resp.ConnectionProfile)
	d.Set("ipsec_psk_a", resp.IpsecPskA)
	d.Set("ipsec_psk_b", resp.IpsecPskB)
	d.Set("cancom_networks", resp.CancomNetworks)
	d.Set("customer_networks", resp.CustomerNetworks)

	return diags
}

// --------------------------Update Connection------------------------------
func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("cmsmgw")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	connectionId := d.Get("id").(string)

	temp1 := d.Get("cancom_networks").([]interface{})
	var cancomNetworks []string
	for _, temp1 := range temp1 {
		cancomNetworks = append(cancomNetworks, temp1.(string))
	}

	temp2 := d.Get("customer_networks").([]interface{})
	var customerNetworks []string
	for _, temp2 := range temp2 {
		customerNetworks = append(customerNetworks, temp2.(string))
	}

	connection := &client_cmsmgw.ConnectionUpdateRequest{
		//MgwId: 					d.Get("mgw_id").(string),
		NameTag:               d.Get("name_tag").(string),
		CustomerPrimaryGwIp:   d.Get("customer_primary_gw_ip").(string),
		CustomerSecondaryGwIp: d.Get("customer_secondary_gw_ip").(string),
		ConnectionProfile:     d.Get("connection_profile").(string),
		IpsecPskA:             d.Get("ipsec_psk_a").(string),
		IpsecPskB:             d.Get("ipsec_psk_b").(string),
		CancomNetworks:        cancomNetworks,
		CustomerNetworks:      customerNetworks,
	}

	resp, err := (*client_cmsmgw.Client)(c).UpdateConnection(connectionId, connection)

	if err != nil {
		return diag.FromErr(err)
	}
	connectionId = resp.ID
	// this is to give dynamo time to get consistant
	//time.Sleep(60 * time.Second)
	resp.Status = "PENDING_DEPLOYMENT"
	err_count := 1
	i := 1
	for resp.Status != "DEPLOYED" {
		time.Sleep(15 * time.Second)
		resp_new, err := (*client_cmsmgw.Client)(c).GetConnection(connectionId)
		tflog.Info(ctx, "Updating Gateway, Iteration: "+strconv.Itoa(i)+" Status: "+resp.Status)
		if err != nil {
			err_count++
			if err_count > 6 {
				return diag.FromErr(err)
			}
			tflog.Info(ctx, "Updating Connection, skipping error, Iteration: "+strconv.Itoa(i)+" Status: "+resp.Status)
			time.Sleep(30 * time.Second)
		} else {
			resp = resp_new
		}
		tflog.Info(ctx, "Updating Connection, Iteration: "+strconv.Itoa(i)+" Status: "+resp.Status)
		// for update, the state is UPDATE_IN_PROGRESS
		if (resp.Status != "UPDATE_IN_PROGRESS") && (resp.Status != "DEPLOYED") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creationg mgw: %s", 404, resp.Status)
			return diag.FromErr(err2)
		}
		i++
		if i >= 40 {
			resp.Status = "UNKNOWN_ERROR"
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}

	d.SetId(resp.ID)

	resourceConnectionRead(ctx, d, m)

	return diags
}

// --------------------------Delete Connection------------------------------
func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("cmsmgw")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//mgwId := d.Id()
	connectionId := d.Get("id").(string)

	err := (*client_cmsmgw.Client)(c).DeleteConnection(connectionId)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err2 := (*client_cmsmgw.Client)(c).GetConnection(connectionId)
	if err2 != nil {
		resp.Status = "UNKNOWN_ERROR"
	}
	i := 1
	for resp.Status == "DEPLOYED" || resp.Status == "DELETE_IN_PROGRESS" {
		time.Sleep(15 * time.Second)
		resp, err2 = (*client_cmsmgw.Client)(c).GetConnection(connectionId)
		if err2 != nil {
			tflog.Info(ctx, "Finished Deleting connection, result: "+err2.Error())
			//diags = append(diags, diag.Diagnostic{
			//	Severity: diag.Warning,
			//	Summary:  "Finished Deleting Translation",
			//	Detail:   err2.Error(),
			//})
			d.SetId("")
			return diags
		} else {
			tflog.Info(ctx, "Deleting Connection, Iteration: "+strconv.Itoa(i)+" Status: "+resp.Status)
		}
		i++
		if i >= 40 {
			resp.Status = "UNKNOWN_ERROR"
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Timeout while waiting for deletion to complete",
				Detail:   "Check the state of resource intended to be deleted",
			})
			return diags
		}
	}

	d.SetId("")
	return diags
}
