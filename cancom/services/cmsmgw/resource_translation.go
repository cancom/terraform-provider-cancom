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

func resourceTranslation() *schema.Resource {
	return &schema.Resource{
		Description:   "Managed Gateway --- Represents a translation resource.",
		CreateContext: resourceTranslationCreate,
		ReadContext:   resourceTranslationRead,
		UpdateContext: resourceTranslationUpdate,
		DeleteContext: resourceTranslationDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mgw_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name_tag": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"customer_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spark_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTranslationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	translation := &client_cmsmgw.TranslationCreateRequest{
		CustomerIp: d.Get("customer_ip").(string),
		MgwId:      d.Get("mgw_id").(string),
		NameTag:    d.Get("name_tag").(string),
		DnsZone:    d.Get("dns_zone").(string),
	}

	//resp, err := c.CreateMgw(CarOrder{
	resp, err := (*client_cmsmgw.Client)(c).CreateTranslation(translation)
	if err != nil {
		return diag.FromErr(err)
	}
	translationId := resp.ID

	//resp.Status = "PENDING"
	i := 1
	for resp.DeploymentState != "DEPLOYED" {
		time.Sleep(15 * time.Second)
		resp, err = (*client_cmsmgw.Client)(c).GetTranslation(translationId)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "Creating Translation, Iteration: "+strconv.Itoa(i)+" Status: "+resp.DeploymentState)
		if (resp.DeploymentState != "PENDING_DEPLOYMENT") && (resp.DeploymentState != "DEPLOYED") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creating translation: %s", 404, resp.DeploymentState)
			return diag.FromErr(err2)
		}

		i++
		if i >= 40 {
			resp.DeploymentState = "UNKNOWN_ERROR" // was ist Status, evtl. nicht richtig
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}

	d.SetId(resp.ID)

	resourceTranslationRead(ctx, d, m)

	return diags
}

func resourceTranslationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	translationId := d.Id()

	//order, err := c.GetOrder(orderId)
	resp, err := (*client_cmsmgw.Client)(c).GetTranslation(translationId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)
	d.Set("customer_ip", resp.CustomerIp)
	d.Set("mgw_id", resp.MgwId)
	d.Set("spark_ip", resp.SparkIp)
	d.Set("name_tag", resp.NameTag)
	d.Set("dns_zone", resp.DnsZone)

	return diags
}

func resourceTranslationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	translationId := d.Get("id").(string)

	translation := &client_cmsmgw.TranslationUpdateRequest{
		CustomerIp: d.Get("customer_ip").(string),
		MgwId:      d.Get("mgw_id").(string),
		NameTag:    d.Get("name_tag").(string),
		DnsZone:    d.Get("dns_zone").(string),
	}

	resp, err := (*client_cmsmgw.Client)(c).UpdateTranslation(translationId, translation)
	if err != nil {
		return diag.FromErr(err)
	}
	translationId = resp.ID
	i := 1
	for resp.DeploymentState != "DEPLOYED" {
		time.Sleep(15 * time.Second)
		resp, err = (*client_cmsmgw.Client)(c).GetTranslation(translationId)
		if err != nil {
			return diag.FromErr(err)
		}
		tflog.Info(ctx, "Updating Translation, Iteration: "+strconv.Itoa(i)+" Status: "+resp.DeploymentState)
		if (resp.DeploymentState != "UPDATE_IN_PROGRESS") && (resp.DeploymentState != "DEPLOYED") {
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Unknown Status creating translation: %s", 404, resp.DeploymentState)
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
			resp.DeploymentState = "UNKNOWN_ERROR" // was ist Status, evtl. nicht richtig
			//signal error and exit, may be we can do something better
			err2 := fmt.Errorf("status: %d, Text: %s", 404, "Timeout")
			return diag.FromErr(err2)
		}
	}

	d.SetId(resp.ID)

	resourceTranslationRead(ctx, d, m)

	return diags
}

func resourceTranslationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("cmsmgw")
	if err != nil {
		return diag.FromErr(err)
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	translationId := d.Get("id").(string)

	err = (*client_cmsmgw.Client)(c).DeleteTranslation(translationId)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err2 := (*client_cmsmgw.Client)(c).GetTranslation(translationId)
	if err2 != nil {
		resp.DeploymentState = "UNKNOWN_ERROR"
	}
	i := 1
	for resp.DeploymentState == "DEPLOYED" || resp.DeploymentState == "DELETE_IN_PROGRESS" {
		time.Sleep(15 * time.Second)
		resp, err2 = (*client_cmsmgw.Client)(c).GetTranslation(translationId)
		if err2 != nil {
			tflog.Info(ctx, "Finished Deleting Translation, result: "+err2.Error())
			//diags = append(diags, diag.Diagnostic{
			//	Severity: diag.Warning,
			//	Summary:  "Finished Deleting Translation",
			//	Detail:   err2.Error(),
			//})
			d.SetId("")
			return diags
		} else {
			tflog.Info(ctx, "Deleting Translation, Iteration: "+strconv.Itoa(i)+" Status: "+resp.DeploymentState)
		}
		i++
		if i >= 40 {
			resp.DeploymentState = "UNKNOWN_ERROR"
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
