package cmsmgw

import (
	"github.com/cancom/terraform-provider-cancom/cancom/services/base"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Provider base.Provider

func (p Provider) Name() string {
	return p.ServiceName
}

func (p Provider) Provider() *schema.Provider {
	return p.ProviderSchema
}

func New() Provider {
	return Provider{
		ServiceName: "cmsmgw",
		ProviderSchema: &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"gateway":     resourceGateway(),
				"connection":  resourceConnection(),
				"translation": resourceTranslation(),
			},
		},
	}
}
