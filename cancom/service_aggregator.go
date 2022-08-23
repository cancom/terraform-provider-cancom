package mcm

import (
	"github.com/cancom/terraform-provider-cancom/cancom/services/base"
	carrental "github.com/cancom/terraform-provider-cancom/cancom/services/car-rental"
	"github.com/cancom/terraform-provider-cancom/cancom/services/dns"
	"github.com/cancom/terraform-provider-cancom/cancom/services/iam"
	sslmonitoring "github.com/cancom/terraform-provider-cancom/cancom/services/ssl-monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getAllProviders() []base.ProviderI {
	return []base.ProviderI{
		sslmonitoring.New(),
		carrental.New(),
		dns.New(),
		iam.New(),
	}
}

func aggregateResources() map[string]*schema.Resource {
	p := getAllProviders()
	resources := make(map[string]*schema.Resource)
	for _, p := range p {
		for k, v := range p.Provider().ResourcesMap {
			k = "cancom_" + p.Name() + "_" + k
			resources[k] = v
		}
	}
	return resources
}
