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

	return diags
}

func resourceBucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("object-storage")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	bucketCreateRequest := client_object_storage.Bucket{
		BucketName:        d.Get("bucket_name").(string),
		AvailabilityClass: d.Get("availability_class").(string),
	}

	resp, err := (*client_object_storage.Client)(c).CreateBucket(&bucketCreateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.BucketName)
	d.Set("bucket_name", resp.BucketName)
	d.Set("availability_class", resp.AvailabilityClass)

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
