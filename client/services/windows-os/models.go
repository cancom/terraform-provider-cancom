package client_windowsos

type WindowsOS_Deplyoment struct {
	Status                 int                `json:"status,omitempty"`
	Id                     string             `json:"deploymentId,omitempty"`
	CoustomerEnvironmentID string             `json:"coustomerEnvironmentID"`
	Computer               WindowsOS_Computer `json:"computer"`
}
type WindowsOS_Computer struct {
	Computername        string   `json:"computername,omitempty"`
	Guid                string   `json:"guid,omitempty"`
	Services            []string `json:"services"`
	CustomerID          string   `json:"customerID,omitempty"`
	MaintenanceWindowId []string `json:"maintenanceWindowId"`
	Role                string   `json:"role"`
}
