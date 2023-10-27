package cancom

import (
	"github.com/cancom/terraform-provider-cancom/cancom/services/base"
	"github.com/cancom/terraform-provider-cancom/cancom/services/cmsmgw"
	"github.com/cancom/terraform-provider-cancom/cancom/services/dns"
	"github.com/cancom/terraform-provider-cancom/cancom/services/iam"
	"github.com/cancom/terraform-provider-cancom/cancom/services/ipam"
	sslmonitoring "github.com/cancom/terraform-provider-cancom/cancom/services/ssl-monitoring"
	windowsos "github.com/cancom/terraform-provider-cancom/cancom/services/windows-os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getAllProviders() []base.ProviderI {
	return []base.ProviderI{
		sslmonitoring.New(),
		dns.New(),
		iam.New(),
		ipam.New(),
		cmsmgw.New(),
		windowsos.New(),
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
