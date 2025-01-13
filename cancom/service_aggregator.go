package cancom

import (
	"github.com/cancom/terraform-provider-cancom/cancom/services/base"
	"github.com/cancom/terraform-provider-cancom/cancom/services/cmsmgw"
	"github.com/cancom/terraform-provider-cancom/cancom/services/dns"
	dynamiccloud "github.com/cancom/terraform-provider-cancom/cancom/services/dynamic-cloud"
	"github.com/cancom/terraform-provider-cancom/cancom/services/iam"
	"github.com/cancom/terraform-provider-cancom/cancom/services/ipam"
	"github.com/cancom/terraform-provider-cancom/cancom/services/s3"
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
		dynamiccloud.New(),
		s3.New(),
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

func aggregateData() map[string]*schema.Resource {
	p := getAllProviders()
	datas := make(map[string]*schema.Resource)
	for _, p := range p {
		for k, v := range p.Provider().DataSourcesMap {
			k = "cancom_" + p.Name() + "_" + k
			datas[k] = v
		}
	}
	return datas
}
