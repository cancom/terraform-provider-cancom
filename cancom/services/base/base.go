package base

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type Provider struct {
	ServiceName    string
	ProviderSchema *schema.Provider
}

type ProviderI interface {
	Name() string
	Provider() *schema.Provider
}
