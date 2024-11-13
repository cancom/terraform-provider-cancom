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

type OpenStackProjectQuotas struct {
	BackupGigabytes    int `json:"backupGigabytes"`
	Backups            int `json:"backups"`
	Cores              int `json:"cores"`
	Floatingip         int `json:"floatingip"`
	Gigabytes          int `json:"gigabytes"`
	HealthMonitor      int `json:"healthMonitor"`
	Listener           int `json:"listener"`
	Loadbalancer       int `json:"loadbalancer"`
	Member             int `json:"member"`
	Pool               int `json:"pool"`
	Instances          int `json:"instances"`
	KeyPairs           int `json:"keyPairs"`
	MetadataItems      int `json:"metadataItems"`
	Network            int `json:"network"`
	PerVolumeGigabytes int `json:"perVolumeGigabytes"`
	Port               int `json:"port"`
	Ram                int `json:"ram"`
	RbacPolicy         int `json:"rbacPolicy"`
	Router             int `json:"router"`
	SecurityGroup      int `json:"securityGroup"`
	SecurityGroupRule  int `json:"securityGroupRule"`
	ServerGroupMembers int `json:"serverGroupMembers"`
	ServerGroups       int `json:"serverGroups"`
	Snapshots          int `json:"snapshots"`
	Subnet             int `json:"subnet"`
	Subnetpool         int `json:"subnetpool"`
	Volumes            int `json:"volumes"`
}

type VpcProjectSpec struct {
	CreatedBy        string                 `json:"createdBy"`
	ProjectComment   string                 `json:"projectComment"`
	OpenstackUuid    string                 `json:"openstackUuid"`
	Limits           OpenStackProjectQuotas `json:"limits"`
	ProjectUsers     []string               `json:"projectUsers"`
	ManagedByService string                 `json:"managedByService,omitempty"`
}

type VpcProjectStatus struct {
	Conditions []Condition            `json:"conditions"`
	Phase      string                 `json:"phase"`
	UpdatedBy  string                 `json:"updatedBy"`
	Usage      OpenStackProjectQuotas `json:"usage"`
	Message    string                 `json:"message,omitempty"`
	Reason     string                 `json:"reason,omitempty"`
}

type VpcProject struct {
	ApiVersion string           `json:"string"`
	Metadata   Metadata         `json:"metadata"`
	Spec       VpcProjectSpec   `json:"spec"`
	Status     VpcProjectStatus `json:"status"`
}

type VpcProjectCreateRequest struct {
	Name           string `json:"name"`
	ProjectComment string `json:"projectComment,omitempty"`
}
