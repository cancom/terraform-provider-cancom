package client_sslmonitoring

type SslMonitor struct {
	IsManagedByCancom    bool   `json:"isManagedByCancom"`
	SslScanEnabled       bool   `json:"sslScanEnabled"`
	Port                 int    `json:"port"`
	CreatedAt            int    `json:"createdAt"`
	LastUpdatedAt        int    `json:"lastUpdatedAt"`
	NotAfter             int    `json:"notAfter"`
	DomainName           string `json:"domainName"`
	Tenant               string `json:"tenant"`
	Comment              string `json:"comment"`
	MinimumGrade         string `json:"minimumGrade"`
	ContactEmailCancom   string `json:"contactEmailCancom"`
	ContactEmailCustomer string `json:"contactEmailCustomer"`
	Protocol             string `json:"protocol"`
	ID                   string `json:"id"`
	State                string `json:"state"`
	SslGrade             string `json:"sslGrade"`
	CreatedBy            string `json:"createdBy"`
	LastUpdatedBy        string `json:"lastUpdatedBy"`
}

type SslMonitorCreateRequest struct {
	SslScanEnabled       bool   `json:"sslScanEnabled"`
	IsManagedByCancom    bool   `json:"isManagedByCancom"`
	Port                 int    `json:"port"`
	DomainName           string `json:"domainName"`
	Tenant               string `json:"tenant"`
	Comment              string `json:"comment"`
	MinimumGrade         string `json:"minimumGrade"`
	ContactEmailCancom   string `json:"contactEmailCancom"`
	ContactEmailCustomer string `json:"contactEmailCustomer"`
	Protocol             string `json:"protocol"`
}
