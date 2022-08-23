package client_customerdb

type Translation struct {
	SourceServiceName   string `json:"sourceServiceName"`
	TargetServiceName   string `json:"targetServiceName"`
	ServiceID           string `json:"serviceId"`
	TranslatedServiceID string `json:"translatedServiceId"`
}

type Customer struct {
	ShortNameID                 string `json:"SHORT_NAME_ID"`
	SnowBusinessUnitDescription string `json:"SNOW_BUSINESS_UNIT_DESCRIPTION"`
	SnowSysID                   string `json:"SNOW_SYS_ID"`
	SnowBusinessUnitName        string `json:"SNOW_BUSINESS_UNIT_NAME"`
	SnowBusinessUnitID          string `json:"SNOW_BUSINESS_UNIT_ID"`
	McmAccountID                string `json:"MCM_ACCOUNT_ID"`
	SnowNumberID                string `json:"SNOW_NUMBER_ID"`
	UReportLanguage             string `json:"u_report_language"`
	Name                        string `json:"name"`
	IsOnboarded                 string `json:"is_onboarded"`
}
