package cancom

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"
	client_serviceregistry "github.com/cancom/terraform-provider-cancom/client/services/service-registry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func Provider() *schema.Provider {

	ar := aggregateResources()

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CANCOM_TOKEN", nil),
			},
			"registry_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://service-registry.portal.cancom.io/v1",
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
			},
		},
		ResourcesMap:         ar,
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	endpoint := d.Get("registry_endpoint").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if token == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create client",
			Detail:   "Token is required",
		})
		return nil, diags
	}

	c, err := client.NewClient(&endpoint, &token)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	services, err := (*client_serviceregistry.Client)(c).GetAllServices()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get services",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	for _, service := range services {
		c.ServiceURLs[service.ServiceName] = service.ServiceEndpoint.Backend
	}

	return c, diags

}
