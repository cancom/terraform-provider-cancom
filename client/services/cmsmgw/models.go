package client_cmsmgw

type Gateway struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GatewayCreateRequest struct {
	Name string `json:"name"`
}
