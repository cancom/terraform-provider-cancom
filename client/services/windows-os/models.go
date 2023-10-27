package client_windowsos

type WindowsOS_Deplyoment struct {
	Status                 int                `json:"status"`
	Id                     string             `json:"deploymentId"`
	CoustomerEnvironmentID string             `json:"coustomerEnvironmentID"`
	Computer               WindowsOS_Computer `json:"computer"`
}
type WindowsOS_Computer struct {
	Computername        string   `json:"computername"`
	Guid                string   `json:"guid"`
	Services            []string `json:"services"`
	CustomerID          string   `json:"customerID"`
	MaintenanceWindowId []string `json:"maintenanceWindowId"`
	Role                string   `json:"role"`
}
