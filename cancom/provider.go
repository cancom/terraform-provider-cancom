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
				Description: "API Token retrieved from [https://portal.cancom.io](https://portal.cancom.io)",
				DefaultFunc: schema.EnvDefaultFunc("CANCOM_TOKEN", nil),
			},
			"service_registry": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service Registry to use for endpoint discovery",
				DefaultFunc: schema.EnvDefaultFunc("CANCOM_SERVICE_REGISTRY", "https://service-registry.portal.cancom.io"),
			},
		},
		ResourcesMap:         ar,
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	service_registry := d.Get("service_registry").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if token != "" {
		c, err := client.NewClient(&service_registry, &token)
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
				Summary:  "Failed to get services " + service_registry,
				Detail:   err.Error(),
			})
			return nil, diags
		}

		for _, service := range services {
			c.ServiceURLs[service.ServiceName] = service.ServiceEndpoint.Backend
		}

		return c, diags
	}

	c, err := client.NewClient(nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
