package client_ipam

type Host struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
}

type HostCreateRequest struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
}

type HostUpdateRequest struct {
	ID          string `json:"crn,omitempty"`
	Operation   string `json:"operation"`
	Address     string `json:"address,omitempty"`
	Qualifier   string `json:"qualifier"`
	NetworkCrn  string `json:"networkCrn"`
	NameTag     string `json:"nameTag"`
	Description string `json:"description"`
}

type HostDeleteResponse struct {
	ID      string `json:"crn,omitempty"`
	Message string `json:"message"`
}

type Instance struct {
	ID              string `json:"crn,omitempty"`
	Description     string `json:"description"`
	NameTag         string `json:"nameTag"`
	ManagedBy       string `json:"managedBy"`
	ReleaseWaitTime int    `json:"releaseWaitTime"`
	CreatedAt       int    `json:"createdAt"`
	UpdatedAt       int    `json:"updatedAt"`
}

type InstanceCreateRequest struct {
	NameTag         string `json:"nameTag"`
	ManagedBy       string `json:"managedBy"`
	Description     string `json:"description"`
	ReleaseWaitTime int    `json:"releaseWaitTime"`
}

type InstanceUpdateRequest struct {
	NameTag         string `json:"nameTag"`
	ManagedBy       string `json:"managedBy"`
	ID              string `json:"crn,omitempty"`
	Description     string `json:"description"`
	ReleaseWaitTime int    `json:"releaseWaitTime"`
}

type InstanceDeleteResponse struct {
	ID      string `json:"crn,omitempty"`
	Message string `json:"message"`
}

type Supernet struct {
	ID           string `json:"crn,omitempty"`
	InstanceId   string `json:"parent"`
	Description  string `json:"description"`
	NameTag      string `json:"nameTag"`
	SupernetCidr string `json:"prefixStr"`
	CreatedAt    int    `json:"createdAt"`
	UpdatedAt    int    `json:"updatedAt"`
}

type SupernetCreateRequest struct {
	ID           string `json:"crn,omitempty"`
	InstanceId   string `json:"parent"`
	Description  string `json:"description"`
	NameTag      string `json:"nameTag"`
	SupernetCidr string `json:"prefixStr"`
}

type SupernetUpdateRequest struct {
	ID           string `json:"crn,omitempty"`
	InstanceId   string `json:"parent"`
	Description  string `json:"description"`
	NameTag      string `json:"nameTag"`
	SupernetCidr string `json:"prefixStr"`
}

type SupernetDeleteResponse struct {
	ID      string `json:"crn,omitempty"`
	Message string `json:"message"`
}

type Network struct {
	ID          string `json:"crn,omitempty"`
	SupernetId  string `json:"parent"`
	Description string `json:"description"`
	NameTag     string `json:"nameTag"`
	PrefixStr   string `json:"prefixStr"`
	HostAssign  bool   `json:"hostAssign"`
	CreatedAt   int    `json:"createdAt"`
	UpdatedAt   int    `json:"updatedAt"`
	//HostAssign   bool	`json:hostAssign`
}

type NetworkCreateRequest struct {
	ID          string `json:"crn,omitempty"`
	SupernetId  string `json:"supernetCrn"`
	Request     string `json:"request"`
	Description string `json:"description"`
	NameTag     string `json:"nameTag"`
	//HostAssign   bool	`json:hostAssign`
}

type NetworkUpdateRequest struct {
	ID          string `json:"crn,omitempty"`
	Description string `json:"description"`
	NameTag     string `json:"nameTag"`
	HostAssign  bool   `json:"hostAssign"`
}

type NetworkDeleteResponse struct {
	ID      string `json:"crn,omitempty"`
	Message string `json:"message"`
}
