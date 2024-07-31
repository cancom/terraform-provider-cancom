package cancom

import (
	"context"

	"github.com/cancom/terraform-provider-cancom/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func Provider() *schema.Provider {

	ar := aggregateResources()

	dr := aggregateData()

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
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role to assume with the provided token. Resources are created with this role instead of the original principal",
			},
		},
		ResourcesMap:         ar,
		DataSourcesMap:       dr,
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	service_registry := d.Get("service_registry").(string)
	// role := d.Get("role").(string)

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

		// if role != "" {
		// 	c.HostURL = c.ServiceURLs["iam"]

		// 	token, err := (*client_iam.Client)(c).AssumeRole(&client_iam.AssumeRoleRequest{
		// 		Role: role,
		// 	})

		// 	if err != nil {
		// 		diags = append(diags, diag.Diagnostic{
		// 			Severity: diag.Error,
		// 			Summary:  fmt.Sprintf("Could not assume role %s", role),
		// 			Detail:   err.Error(),
		// 		})
		// 		return nil, diags
		// 	}

		// 	c.Token = token.Jwt
		// }

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
