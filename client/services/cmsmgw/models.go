package client_cmsmgw

type Gateway struct {
	ID                  string `json:"mgwId"`
	Message             string `json:"message,omitempty"`
	NameTag             string `json:"nameTag,omitempty"`
	Customer            string `json:"customer,omitempty"`
	State               string `json:"state"`
	CancomPrimaryGwIp   string `json:"primaryGwPublicIp,omitempty"`
	CancomSecondaryGwIp string `json:"secondaryGwPublicIp,omitempty"`
	BastionNetwork      string `json:"remoteBastionNetwork,omitempty"`
	MgwSize             string `json:"mgwSize,omitempty"`
	Tag                 string `json:"tag,omitempty"`
	NatTranslation      bool   `json:"natTranslation"`
	BastionLiteLinux    bool   `json:"bastionLiteLinux"`
	BastionLiteWindows  bool   `json:"bastionLiteWindows"`
	NatNetwork          string `json:"natNetwork,omitempty"`
}

type GatewayCreateRequest struct {
	ID                 string `json:"mgwId"`
	NameTag            string `json:"nameTag,omitempty"`
	Customer           string `json:"customer,omitempty"`
	State              string `json:"state"`
	MgwSize            string `json:"mgwSize,omitempty"`
	Tag                string `json:"tag,omitempty"`
	NatTranslation     bool   `json:"natTranslation"`
	BastionLiteLinux   bool   `json:"bastionLiteLinux"`
	BastionLiteWindows bool   `json:"bastionLiteWindows"`
}

type GatewayUpdateRequest struct {
	ID                 string `json:"mgwId"`
	NameTag            string `json:"nameTag,omitempty"`
	Customer           string `json:"customer,omitempty"`
	States             string `json:"state"`
	MgwSize            string `json:"mgwSize,omitempty"`
	Tag                string `json:"tag,omitempty"`
	NatTranslation     bool   `json:"natTranslation"`
	BastionLiteLinux   bool   `json:"bastionLiteLinux"`
	BastionLiteWindows bool   `json:"bastionLiteWindows"`
}

type Translation struct {
	ID              string `json:"translationId,omitempty"`
	MgwId           string `json:"mgwId,omitempty"`
	NameTag         string `json:"nameTag,omitempty"`
	SparkIp         string `json:"sparkIp,omitempty"`
	CustomerIp      string `json:"customerIp,omitempty"`
	DeploymentState string `json:"deploymentState,omitempty"`
	DnsZone         string `json:"dnsZone,omitempty"`
}

type TranslationCreateRequest struct {
	MgwId      string `json:"mgwId,omitempty"`
	NameTag    string `json:"nameTag,omitempty"`
	SparkIp    string `json:"sparkIp,omitempty"`
	CustomerIp string `json:"customerIp,omitempty"`
	DnsZone    string `json:"dnsZone,omitempty"`
}

type TranslationUpdateRequest struct {
	MgwId      string `json:"mgwId,omitempty"`
	NameTag    string `json:"nameTag,omitempty"`
	SparkIp    string `json:"sparkIp,omitempty"`
	CustomerIp string `json:"customerIp,omitempty"`
	DnsZone    string `json:"dnsZone,omitempty"`
}

type Connection struct {
	ID                    string   `json:"connectionId,omitempty"`
	MgwId                 string   `json:"mgwId,omitempty"`
	NameTag               string   `json:"nameTag,omitempty"`
	CustomerPrimaryGwIp   string   `json:"customerPrimaryGwIp,omitempty"`
	CustomerSecondaryGwIp string   `json:"customerSecondaryGwIp"`
	Status                string   `json:"deploymentState,omitempty"`
	ConnectionProfile     string   `json:"connectionProfile,omitempty"`
	IpsecPskA             string   `json:"ipsecPskA,omitempty"`
	IpsecPskB             string   `json:"ipsecPskB,omitempty"`
	CancomNetworks        []string `json:"cancomNetworks"`
	CustomerNetworks      []string `json:"customerNetworks"`
}

type ConnectionCreateRequest struct {
	ID                    string   `json:"connectionId,omitempty"`
	MgwId                 string   `json:"mgwId,omitempty"`
	NameTag               string   `json:"nameTag,omitempty"`
	CustomerPrimaryGwIp   string   `json:"customerPrimaryGwIp,omitempty"`
	CustomerSecondaryGwIp string   `json:"customerSecondaryGwIp"`
	ConnectionProfile     string   `json:"connectionProfile,omitempty"`
	IpsecPskA             string   `json:"ipsecPskA,omitempty"`
	IpsecPskB             string   `json:"ipsecPskB,omitempty"`
	CancomNetworks        []string `json:"cancomNetworks"`
	CustomerNetworks      []string `json:"customerNetworks"`
}

type ConnectionUpdateRequest struct {
	ID string `json:"connectionId,omitempty"`
	//MgwId                 string   `json:"mgwId,omitempty"`
	NameTag               string   `json:"nameTag,omitempty"`
	CustomerPrimaryGwIp   string   `json:"customerPrimaryGwIp,omitempty"`
	CustomerSecondaryGwIp string   `json:"customerSecondaryGwIp"`
	ConnectionProfile     string   `json:"connectionProfile,omitempty"`
	IpsecPskA             string   `json:"ipsecPskA,omitempty"`
	IpsecPskB             string   `json:"ipsecPskB,omitempty"`
	CancomNetworks        []string `json:"cancomNetworks"`
	CustomerNetworks      []string `json:"customerNetworks"`
}
