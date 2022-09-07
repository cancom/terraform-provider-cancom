package client_serviceregistry

type ServiceEndpoint struct {
	Frontend string `json:"frontend"`
	Backend  string `json:"backend"`
}

type Service struct {
	ServiceName             string          `json:"serviceName"`
	DisplayName             string          `json:"displayName"`
	Route                   string          `json:"route"`
	ServiceDocumentation    string          `json:"serviceDocumentation"`
	RelatedIamPermissions   []string        `json:"relatedIamPermissions"`
	RelatedServices         []string        `json:"relatedServices"`
	RequiredServiceAccounts []string        `json:"requiredServiceAccounts"`
	Keywords                []string        `json:"keywords"`
	ServiceEndpoint         ServiceEndpoint `json:"serviceEndpoint"`
}

type CreateServiceBody struct {
	OverwriteService        bool            `json:"overwriteService"`
	DisplayName             string          `json:"displayName"`
	ServiceDocumentation    string          `json:"serviceDocumentation"`
	Route                   string          `json:"route"`
	RelatedIamPermissions   []string        `json:"relatedIamPermissions"`
	RelatedServices         []string        `json:"relatedServices"`
	RequiredServiceAccounts []string        `json:"requiredServiceAccounts"`
	Keywords                []string        `json:"keywords"`
	ServiceEndpoint         ServiceEndpoint `json:"serviceEndpoint"`
}
