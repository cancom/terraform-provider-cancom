package dns

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dns "github.com/cancom/terraform-provider-cancom/client/services/dns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the record",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the record (i.e. A, CNAME, ...)",
				Required:    true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_change_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["domdns"]

	var diags diag.Diagnostics

	record := &client_dns.RecordCreateRequest{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Content:  d.Get("content").(string),
		TTL:      d.Get("ttl").(int),
		ZoneName: d.Get("zone_name").(string),
	}

	resp, err := (*client_dns.Client)(c).CreateRecord(record)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.Set("last_change_date", resp.LastChangeDate)
	d.Set("zone_id", resp.ZoneID)
	d.Set("comments", resp.Comments)

	d.SetId(id)

	resourceRecordRead(ctx, d, m)

	return diags

}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["domdns"]

	var diags diag.Diagnostics

	id := d.Id()
	zoneName := d.Get("zone_name").(string)

	resp, err := (*client_dns.Client)(c).GetRecord(id, zoneName)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("content", resp.Content)
	d.Set("ttl", resp.TTL)
	d.Set("zone_name", resp.ZoneName)
	d.Set("zone_id", resp.ZoneID)
	d.Set("id", resp.ID)
	d.Set("comments", resp.Comments)
	d.Set("last_change_date", resp.LastChangeDate)

	return diags
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["domdns"]

	var diags diag.Diagnostics

	id := d.Id()

	record := &client_dns.RecordUpdateRequest{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Content:  d.Get("content").(string),
		TTL:      d.Get("ttl").(int),
		ZoneName: d.Get("zone_name").(string),
		ZoneID:   d.Get("zone_id").(string),
	}

	_, err := (*client_dns.Client)(c).UpdateRecord(id, record)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceRecordRead(ctx, d, m)

	return diags
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["domdns"]

	var diags diag.Diagnostics

	id := d.Id()

	err := (*client_dns.Client)(c).DeleteRecord(id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
