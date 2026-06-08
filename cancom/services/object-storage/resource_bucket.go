package object_storage

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_object_storage "github.com/cancom/terraform-provider-cancom/client/services/object-storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBucket() *schema.Resource {
	return &schema.Resource{
		Description:   "Object Storage --- Create and manage S3 buckets. Bucket names must be globally unique. ",
		CreateContext: resourceBucketCreate,
		ReadContext:   resourceBucketRead,
		UpdateContext: resourceBucketUpdate,
		DeleteContext: resourceBucketDelete,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the bucket. The name must be unique GLOBALLY.",
			},
			"availability_class": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the availability class. One of singleDc or multiDc.",
				ValidateFunc: validation.StringInSlice([]string{"singleDc", "multiDc"}, false),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description what the bucket is used for",
			},
			"ip_whitelist": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "A list of CIDRs that should be able to access the bucket",
			},
		},
	}
}

func resourceBucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	resp, err := (*client_object_storage.Client)(c).GetBucket(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("bucket_name", resp.BucketName)
	d.Set("availability_class", resp.AvailabilityClass)
	d.Set("description", resp.Description)
	d.Set("ip_whitelist", resp.IpWhitelist)

	return diags
}

func resourceBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	// 1. Grab the raw slice of interfaces
	rawWhitelist := d.Get("ip_whitelist").([]interface{})
	ipWhitelist := make([]string, len(rawWhitelist))

	// 2. Safely assert each element to a string
	for i, val := range rawWhitelist {
		ipWhitelist[i] = val.(string)
	}

	bucketCreateRequest := client_object_storage.Bucket{
		BucketName:        d.Get("bucket_name").(string),
		AvailabilityClass: d.Get("availability_class").(string),
		Description:       d.Get("description").(string),
		IpWhitelist:       ipWhitelist,
	}

	resp, err := (*client_object_storage.Client)(c).CreateBucket(&bucketCreateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.BucketName)
	d.Set("bucket_name", resp.BucketName)
	d.Set("availability_class", resp.AvailabilityClass)
	d.Set("description", resp.Description)
	d.Set("ip_whitelist", resp.IpWhitelist)

	return diags
}

func resourceBucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	// 1. Grab the raw slice of interfaces
	rawWhitelist := d.Get("ip_whitelist").([]interface{})
	ipWhitelist := make([]string, len(rawWhitelist))

	// 2. Safely assert each element to a string
	for i, val := range rawWhitelist {
		ipWhitelist[i] = val.(string)
	}

	bucketUpdateRequest := client_object_storage.BucketUpdateRequest{
		Description: d.Get("description").(string),
		IpWhitelist: ipWhitelist,
	}

	resp, err := (*client_object_storage.Client)(c).UpdateBucket(d.Id(), bucketUpdateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", resp.Description)
	d.Set("ip_whitelist", resp.IpWhitelist)

	return diags
}

func resourceBucketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	err = (*client_object_storage.Client)(c).DeleteBucket(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
