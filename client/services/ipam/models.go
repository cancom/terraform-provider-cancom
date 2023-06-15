package client_ipam

type Gateway struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SupernetCreateRequest struct {
	Name string `json:"name"`
}

type Host struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

type HostCreateRequest struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

type HostUpdateRequest struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
	Source      string `json:"source"`
}
