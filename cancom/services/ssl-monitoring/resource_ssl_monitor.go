package sslmonitoring

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_sslmonitoring "github.com/cancom/terraform-provider-cancom/client/services/ssl-monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSslMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSslMonitorCreate,
		ReadContext:   resourceSslMonitorRead,
		UpdateContext: resourceSslMonitorUpdate,
		DeleteContext: resourceSslMonitorDelete,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name to monitor. A domain name can only be monitored once.",
			},
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tenant that created the monitoring item.",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment that provides more details about the item, for example what to do when an alarm is raised",
			},
			"ssl_scan_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Are scans scheduled, or do we not schedule them at this moment?.",
			},
			"minimum_grade": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alert grade. If the score falls below this score, an alarm is raised.",
			},
			"contact_email_cancom": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CANCOM contact person that is responsible to exchange the certificate.",
			},
			"contact_email_customer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notification is send to the contact if an alarm is raised.",
			},
			"is_managed_by_cancom": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "This item is managed by CANCOM, and a CANCOM employee is responible to renew the certificate. Can only be set by CANCOM emplyoees.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol to be monitored, i.e. `http`, `postgresql`, `smtp`, ...",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port to scan.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the monitoring item.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State, for example `finished` or `running`.",
			},
			"not_after": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Information from the certificate, when it is about to expire.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who created the monitoring item.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "When was the monitoring item created.",
			},
			"last_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Who was the last principal to update the item.",
			},
			"last_updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Datetime of the last update.",
			},
			"auto_scan": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Only relevant on resource creation. Automatically starts a scan while creating the monitoring item. *This operation incurs a charge*",
			},
		},
	}
}

func resourceSslMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("ssl-monitoring")

	var diags diag.Diagnostics

	sslMonitorRequest := client_sslmonitoring.SslMonitorCreateRequest{
		DomainName:           d.Get("domain_name").(string),
		Tenant:               d.Get("tenant").(string),
		Comment:              d.Get("comment").(string),
		SslScanEnabled:       d.Get("ssl_scan_enabled").(bool),
		MinimumGrade:         d.Get("minimum_grade").(string),
		ContactEmailCancom:   d.Get("contact_email_cancom").(string),
		ContactEmailCustomer: d.Get("contact_email_customer").(string),
		IsManagedByCancom:    d.Get("is_managed_by_cancom").(bool),
		Protocol:             d.Get("protocol").(string),
		Port:                 d.Get("port").(int),
	}

	resp, err := (*client_sslmonitoring.Client)(c).CreateSslMonitor(&sslMonitorRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.ID)

	if d.Get("auto_scan").(bool) {
		err := (*client_sslmonitoring.Client)(c).StartSslScan(resp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags

}

func resourceSslMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("ssl-monitoring")

	var diags diag.Diagnostics

	sslMonitorRequest := client_sslmonitoring.SslMonitor{
		ID: d.Id(),
	}

	resp, err := (*client_sslmonitoring.Client)(c).GetSslMonitor(sslMonitorRequest.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("domain_name", resp.DomainName)
	d.Set("tenant", resp.Tenant)
	d.Set("comment", resp.Comment)
	d.Set("ssl_scan_enabled", resp.SslScanEnabled)
	d.Set("minimum_grade", resp.MinimumGrade)
	d.Set("contact_email_cancom", resp.ContactEmailCancom)
	d.Set("contact_email_customer", resp.ContactEmailCustomer)
	d.Set("is_managed_by_cancom", resp.IsManagedByCancom)
	d.Set("protocol", resp.Protocol)
	d.Set("port", resp.Port)
	d.Set("id", resp.ID)
	d.Set("state", resp.State)
	d.Set("not_after", resp.NotAfter)
	d.Set("created_by", resp.CreatedBy)
	d.Set("created_at", resp.CreatedAt)
	d.Set("last_updated_by", resp.LastUpdatedBy)
	d.Set("last_updated_at", resp.LastUpdatedAt)

	return diags
}

func resourceSslMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("ssl-monitoring")

	var diags diag.Diagnostics

	sslMonitorRequest := client_sslmonitoring.SslMonitor{
		ID: d.Id(),
	}

	sslMonitorUpdateRequest := client_sslmonitoring.SslMonitorCreateRequest{}

	sslMonitorUpdateRequest.DomainName = d.Get("domain_name").(string)
	sslMonitorUpdateRequest.Tenant = d.Get("tenant").(string)
	sslMonitorUpdateRequest.Comment = d.Get("comment").(string)
	sslMonitorUpdateRequest.SslScanEnabled = d.Get("ssl_scan_enabled").(bool)
	sslMonitorUpdateRequest.MinimumGrade = d.Get("minimum_grade").(string)
	sslMonitorUpdateRequest.ContactEmailCancom = d.Get("contact_email_cancom").(string)
	sslMonitorUpdateRequest.ContactEmailCustomer = d.Get("contact_email_customer").(string)
	sslMonitorUpdateRequest.IsManagedByCancom = d.Get("is_managed_by_cancom").(bool)
	sslMonitorUpdateRequest.Protocol = d.Get("protocol").(string)
	sslMonitorUpdateRequest.Port = d.Get("port").(int)

	resp, err := (*client_sslmonitoring.Client)(c).UpdateSslMonitor(sslMonitorRequest.ID, &sslMonitorUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated_at", resp.LastUpdatedAt)

	d.SetId(resp.ID)

	if d.Get("auto_scan").(bool) {
		err := (*client_sslmonitoring.Client)(c).StartSslScan(resp.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceSslMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, _ := m.(*client.CcpClient).GetService("ssl-monitoring")

	var diags diag.Diagnostics

	sslMonitorRequest := client_sslmonitoring.SslMonitor{
		ID: d.Id(),
	}

	err := (*client_sslmonitoring.Client)(c).DeleteSslMonitor(sslMonitorRequest.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
