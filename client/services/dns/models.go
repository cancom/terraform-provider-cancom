package client_dns

type Zone struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	Type           string `json:"type"`
	Tenant         string `json:"tenant"`
	Email          string `json:"email"`
	SOA            SOA    `json:"soa"`
	DnsSec         string `json:"dnsSec"`
	LastChangeDate string `json:"lastChangeDate"`
}

type SOA struct {
	Refresh     int `json:"refresh"`
	Retry       int `json:"retry"`
	Expire      int `json:"expire"`
	TTL         int `json:"ttl"`
	NegativeTTL int `json:"negative_ttl"`
}

type Record struct {
	ID             string `json:"id"`
	ZoneID         string `json:"zoneId"`
	ZoneName       string `json:"zoneName"`
	Type           string `json:"type"`
	Name           string `json:"name"`
	Content        string `json:"content"`
	TTL            int    `json:"ttl"`
	Comments       string `json:"comments"`
	LastChangeDate string `json:"lastChangeDate"`
	Tenant         string `json:"tenant"`
}

type RecordCreateRequest struct {
	ZoneName string `json:"zoneName"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Tenant   string `json:"tenant"`
	Mode     string `json:"mode"`
}

type RecordUpdateRequest struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	ZoneName string `json:"zoneName"`
	ZoneID   string `json:"zoneId"`
	Tenant   string `json:"tenant"`
	Mode     string `json:"mode"`
}
