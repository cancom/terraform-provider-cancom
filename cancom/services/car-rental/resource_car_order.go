package carrental

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_carrental "github.com/cancom/terraform-provider-cancom/client/services/car-rental"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCarOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCarOrderCreate,
		ReadContext:   resourceCarOrderRead,
		UpdateContext: resourceCarOrderUpdate,
		DeleteContext: resourceCarOrderDelete,
		Schema: map[string]*schema.Schema{
			"order_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vehicle_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hp": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"seats": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mileage": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceCarOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["car-rental"]

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	order := client_carrental.CarOrderCreateRequest{
		OrderName:    d.Get("order_name").(string),
		Type:         d.Get("type").(string),
		VehicleClass: d.Get("vehicle_class").(string),
		HP:           d.Get("hp").(int),
		Seats:        d.Get("seats").(string),
	}

	resp, err := (*client_carrental.Client)(c).CreateOrder(order)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(id)

	resourceCarOrderRead(ctx, d, m)

	return diags
}

func resourceCarOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["car-rental"]

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderId := d.Id()

	order, err := (*client_carrental.Client)(c).GetOrder(orderId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("order_name", order.OrderName)
	d.Set("type", order.Type)
	d.Set("vehicle_class", order.VehicleClass)
	d.Set("hp", order.HP)
	d.Set("seats", order.Seats)
	d.Set("id", order.ID)
	d.Set("mileage", order.Mileage)

	return diags
}

func resourceCarOrderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["car-rental"]

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	order := client_carrental.CarOrderUpdateRequest{
		OrderID:      d.Id(),
		OrderName:    d.Get("order_name").(string),
		Type:         d.Get("type").(string),
		VehicleClass: d.Get("vehicle_class").(string),
		HP:           d.Get("hp").(int),
		Seats:        d.Get("seats").(string),
		Mileage:      d.Get("mileage").(float64),
	}

	resp, err := (*client_carrental.Client)(c).UpdateOrder(order)

	if err != nil {
		return diag.FromErr(err)
	}

	id := resp.ID

	d.SetId(id)

	resourceCarOrderRead(ctx, d, m)

	return diags
}

func resourceCarOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["car-rental"]

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderId := d.Id()

	err := (*client_carrental.Client)(c).DeleteOrder(orderId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
