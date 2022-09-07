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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_scan_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"minimum_grade": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_email_cancom": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_email_customer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_managed_by_cancom": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"not_after": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_scan": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceSslMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ssl-monitoring"]

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
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ssl-monitoring"]

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
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ssl-monitoring"]

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
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["ssl-monitoring"]

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
