package client

// TODO: ssl-report

type CI struct {
	CreatedAt   uint32 `json:"created_at"`
	Heartbeat   uint32 `json:"heartbeat"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateConfigurationItemsRequestContent struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	Details      struct {
		Description string `json:"description"`
	} `json:"details"`
	Tenant string `json:"tenant"`
}

type DeleteConfigurationItemsRequestContent struct {
	ID           string `json:"id"`
	Service      string `json:"service"`
	ResourceType string `json:"resourceType"`
	Tenant       string `json:"tenant"`
}
