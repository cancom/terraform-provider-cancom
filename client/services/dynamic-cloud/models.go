package client_dynamiccloud

type Metadata struct {
	CreationDate string `json:"creationDate"`
	DeletionDate string `json:"deletionDate,omitempty"`
	Name         string `json:"name"`
	Shortid      string `json:"shortid"`
	Tenant       string `json:"tenant"`
}

type Condition struct {
	LastTransitionTime string `json:"lastTransitionTime"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	// A human readable message indicating details about why the project is in this condition.
	Message string `json:"message,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
}

type VpcProjectSpec struct {
	CreatedBy        string         `json:"createdBy"`
	ProjectComment   string         `json:"projectComment"`
	OpenstackUuid    string         `json:"openstackUuid"`
	Limits           map[string]int `json:"limits"`
	ProjectUsers     []string       `json:"projectUsers"`
	ManagedByService string         `json:"managedByService,omitempty"`
}

type VpcProjectStatus struct {
	Conditions []Condition    `json:"conditions"`
	Phase      string         `json:"phase"`
	UpdatedBy  string         `json:"updatedBy"`
	Usage      map[string]int `json:"usage"`
	Message    string         `json:"message,omitempty"`
	Reason     string         `json:"reason,omitempty"`
}

type VpcProject struct {
	ApiVersion string           `json:"string"`
	Metadata   Metadata         `json:"metadata"`
	Spec       VpcProjectSpec   `json:"spec"`
	Status     VpcProjectStatus `json:"status"`
}

type VpcProjectCreateMetadata struct {
	Name string `json:"name"`
}

type VpcProjectCreateSpec struct {
	ProjectComment string `json:"projectComment"`
}

type VpcProjectCreateRequest struct {
	Metadata VpcProjectCreateMetadata `json:"metadata"`
	Spec     VpcProjectCreateSpec     `json:"spec"`
}
