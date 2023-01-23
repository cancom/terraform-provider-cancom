package client_ipam

type Gateway struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SupernetCreateRequest struct {
	Name string `json:"name"`
}
