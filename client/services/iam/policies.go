package client_iam

type PolicyRequest struct {
	Policy  map[string]string `json:"policy"`
	Service string            `json:"service"`
}

type CarrentalPolicy struct {
	CreateCarOrders      string `json:"createCarOrders"`
	DeleteCarOrders      string `json:"deleteCarOrders"`
	DescribeCarOrders    string `json:"describeCarOrders"`
	FastForward          string `json:"fastForward"`
	ListCarOrders        string `json:"listCarOrders"`
	ManageFixedContracts string `json:"managefixedContracts"`
	UpdateCarOrders      string `json:"updateCarOrders"`
}

type SslMonitoringPolicy struct {
	DescribeMonitors     string `json:"describeMonitors"`
	CreateMonitors       string `json:"createMonitors"`
	UpdateMonitors       string `json:"updateMonitors"`
	DeleteMonitors       string `json:"deleteMonitors"`
	InternalMonitorAdmin string `json:"internalMonitorAdmin"`
	InternalReportWorker string `json:"internalReportWorker"`
}

type AzureAppManagementPolicy struct {
	ApplyApplicationManagement    string `json:"applyApplicationManagement"`
	CreateApplicationManagement   string `json:"createApplicationManagement"`
	CreateCustomer                string `json:"createCustomer"`
	CreateEnvironment             string `json:"createEnvironment"`
	DescribeApplicationManagement string `json:"describeApplicationManagement"`
	DescribeCustomer              string `json:"describeCustomer"`
	DescribeEnvironment           string `json:"describeEnvironment"`
	ListCustomers                 string `json:"listCustomers"`
	ListEnvironment               string `json:"listEnvironment"`
	UpdateApplicationManagement   string `json:"updateApplicationManagement"`
	UpdateCustomer                string `json:"updateCustomer"`
	UpdateEnvironment             string `json:"updateEnvironment"`
}

type AzureLandingZonePolicy struct {
	CreateLandingZone string `json:"createLandingZone"`
	DeleteLandingZone string `json:"deleteLandingZone"`
	ReadLandingZone   string `json:"readLandingZone"`
	UpdateLandingZone string `json:"updateLandingZone"`
}

type AzureGovernancePolicy struct {
	ApplyDeployment           string `json:"applyDeployment"`
	CreateCustomer            string `json:"createCustomer"`
	CreateDeployment          string `json:"createDeployment"`
	CreateEnvironment         string `json:"createEnvironment"`
	CreateEnvironmentSchedule string `json:"createEnvironmentSchedule"`
	DeleteEnvironmentSchedule string `json:"deleteEnvironmentSchedule"`
	DescribeCustomer          string `json:"describeCustomer"`
	DescribeDeployment        string `json:"describeDeployment"`
	DescribeEnvironment       string `json:"describeEnvironment"`
	DescribeTerracoreProject  string `json:"describeTerracoreProject"`
	DescribeTerracoreSchedule string `json:"describeTerracoreSchedule"`
	ListCustomers             string `json:"listCustomers"`
	ListDataSync              string `json:"listDataSync"`
	ListEnvironment           string `json:"listEnvironment"`
	UpdateCustomer            string `json:"updateCustomer"`
	UpdateDeployment          string `json:"updateDeployment"`
	UpdateEnvironment         string `json:"updateEnvironment"`
	UpdateEnvironmentSchedule string `json:"updateEnvironmentSchedule"`
	WriteDataSync             string `json:"writeDataSync"`
}

type ManagedGatewayPolicy struct {
	InternalRead     string `json:"internalRead"`
	InternalWrite    string `json:"internalWrite"`
	ReadConnection   string `json:"readConnection"`
	ReadMgw          string `json:"readMgw"`
	ReadSpark        string `json:"readSpark"`
	ReadTranslation  string `json:"readTranslation"`
	WriteConnection  string `json:"writeConnection"`
	WriteMgw         string `json:"writeMgw"`
	WriteSpark       string `json:"writeSpark"`
	WriteTranslation string `json:"writeTranslation"`
}

type ContainerOrchestrationPolicy struct {
	CreateClusters       string `json:"createClusters"`
	CreateClusterUsers   string `json:"createClusterUsers"`
	DeleteClusters       string `json:"deleteClusters"`
	DeleteClusterUsers   string `json:"deleteClusterUsers"`
	DescribeClusters     string `json:"describeClusters"`
	DescribeClusterUsers string `json:"describeClusterUsers"`
	ListClusters         string `json:"listClusters"`
	ListClusterUsers     string `json:"listClusterUsers"`
	ListMyUsers          string `json:"listMyUsers"`
	UpdateClusterUsers   string `json:"updateClusterUsers"`
}

type CustomerDBPolicy struct {
	InternalWriteAccess string `json:"internalWriteAccess"`
	ReadAccess          string `json:"readAccess"`
}

type DnsPolicy struct {
	CreateRecords     string `json:"createRecords"`
	CreateZones       string `json:"createZones"`
	DeleteRecords     string `json:"deleteRecords"`
	DescribeRecords   string `json:"describeRecords"`
	DescribeZones     string `json:"describeZones"`
	InternalIdMapping string `json:"internalIdMapping"`
	ListRecords       string `json:"listRecords"`
	ListZones         string `json:"listZones"`
	UpdateRecords     string `json:"updateRecords"`
}

type DrivePolicy struct {
	ReadAccess              string `json:"readAccess"`
	WriteAccess             string `json:"writeAccess"`
	InternalRootReadAccess  string `json:"internalRootReadAccess"`
	InternalRootWriteAccess string `json:"internalRootWriteAccess"`
}

type EventBridgePolicy struct {
	CreateEventConfig string `json:"createEventConfig"`
	DeleteEventConfig string `json:"deleteEventConfig"`
	ReadEventConfig   string `json:"readEventConfig"`
	UpdateEventConfig string `json:"updateEventConfig"`
	InternalRead      string `json:"internalRead"`
}

type GenericServicePolicy struct {
	CreateGeneric string `json:"createGeneric"`
	DeleteGeneric string `json:"deleteGeneric"`
	ReadGeneric   string `json:"readGeneric"`
	UpdateGeneric string `json:"updateGeneric"`
}

type IamPolicy struct {
	AssumeRole                     string `json:"assumeRole"`
	CreateManagedRoles             string `json:"createManagedRoles"`
	CreateRoles                    string `json:"createRoles"`
	CreateServiceUsers             string `json:"createServiceUsers"`
	CreateSession                  string `json:"createSession"`
	CreateUsers                    string `json:"createUsers"`
	DeleteAllSessionsForUser       string `json:"deleteAllSessionsForUser"`
	DeleteManagedRoles             string `json:"deleteManagedRoles"`
	DeletePolicyDocument           string `json:"deletePolicyDocument"`
	DeleteRoles                    string `json:"deleteRoles"`
	DeleteServiceUsers             string `json:"deleteServiceUsers"`
	DeleteUsers                    string `json:"deleteUsers"`
	DescribePolicyDocument         string `json:"describePolicyDocument"`
	DescribeRoles                  string `json:"describeRoles"`
	DescribeServiceUsers           string `json:"describeServiceUsers"`
	DescribeTrustRelations         string `json:"describeTrustRelations"`
	DescribeTrustRelationsAdmin    string `json:"describeTrustRelationsAdmin"`
	DescribeUsers                  string `json:"describeUsers"`
	InternalAssumeRole             string `json:"internalAssumeRole"`
	InternalBypassUpdatePolicy     string `json:"internalBypassUpdatePolicy"`
	InternalDescribeTrustRelations string `json:"internalDescribeTrustRelations"`
	InternalManageTenants          string `json:"internalManageTenants"`
	InternalRoleConnections        string `json:"internalRoleConnections"`
	InternalUpdateTrustRelation    string `json:"internalUpdateTrustRelation"`
	ListPolicyDocument             string `json:"listPolicyDocument"`
	ListRoles                      string `json:"listRoles"`
	ListServiceUsers               string `json:"listServiceUsers"`
	ListUsers                      string `json:"listUsers"`
	ReadSessions                   string `json:"readSessions"`
	RoleConnections                string `json:"roleConnections"`
	RollSession                    string `json:"rollSession"`
	UpdateManagedRoles             string `json:"updateManagedRoles"`
	UpdatePolicyDocument           string `json:"updatePolicyDocument"`
	UpdateRoles                    string `json:"updateRoles"`
	UpdateServiceUsers             string `json:"updateServiceUsers"`
	UpdateTrustRelation            string `json:"updateTrustRelation"`
	UpdateUsers                    string `json:"updateUsers"`
	UpdateSessions                 string `json:"updateSessions"`
}

type ManagedOSPolicy struct {
	DeleteVmwareInstances   string `json:"deleteVmwareInstances"`
	CreateVmwareInstances   string `json:"createVmwareInstances"`
	UpdateVmwareInstances   string `json:"updateVmwareInstances"`
	ListVmwareInstances     string `json:"listVmwareInstances"`
	DescribeVmwareInstances string `json:"describeVmwareInstances"`
}

type MeteringEnginePolicy struct {
	ReportMetering              string `json:"reportMetering"`
	DescribeMetering            string `json:"describeMetering"`
	WriteBillableObjects        string `json:"writeBillableObjects"`
	InternalReadBillableObjects string `json:"internalReadBillableObjects"`
}

type MiaMuseumPolicy struct {
	ReadAccess string `json:"readAccess"`
}

type ServiceRegistryPolicy struct {
	InternalWriteAccess string `json:"internalWriteAccess"`
	WriteAccess         string `json:"writeAccess"`
	ExplicitReadAccess  string `json:"explicitReadAccess"`
}

type SmsNotificationsPolicy struct {
	ManageGroups      string `json:"manageGroups"`
	ManageContacts    string `json:"manageContacts"`
	SendNotifications string `json:"sendNotifications"`
	ManageTemplates   string `json:"manageTemplates"`
}

type ResourceExplorerPolicy struct {
	Read          string `json:"read"`
	InternalWrite string `json:"internalWrite"`
}

type TerracorePolicy struct {
	UpdateProjectTokens   string `json:"updateProjectTokens"`
	DescribeProject       string `json:"describeProject"`
	DeleteProject         string `json:"deleteProject"`
	UpdateProject         string `json:"updateProject"`
	CreateExecution       string `json:"createExecution"`
	CreateProject         string `json:"createProject"`
	ListProjects          string `json:"listProjects"`
	ApplyExecution        string `json:"applyExecution"`
	DescribeExecutionPlan string `json:"describeExecutionPlan"`
	DescribeExecution     string `json:"describeExecution"`
	RestoreProject        string `json:"restoreProject"`
	DeleteSchedule        string `json:"deleteSchedule"`
	CreateSchedule        string `json:"createSchedule"`
	UpdateSchedule        string `json:"updateSchedule"`
	DescribeSchedule      string `json:"describeSchedule"`
	ListSchedules         string `json:"listSchedules"`
}
