package iam

import (
	"context"
	"reflect"
	"unicode"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service name of the ",
			},
			"policy": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Can be set to `reader`, or `full-access`. Applies either reader or administrative access to the service.",
						},
						"custom": {
							Type:        schema.TypeMap,
							Optional:    true,
							ForceNew:    true,
							Description: "Provide a customized policy. If `custom` is set, it is preferred over `profile`",
						},
					},
				},
			},
			"principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The principal that the policy should be applied to. Can be applied to Users, ServiceUsers, and Roles.",
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	var diags diag.Diagnostics

	policy := d.Get("policy").([]interface{})
	profile := policy[0].(map[string]interface{})["profile"].(string)
	custom := policy[0].(map[string]interface{})["custom"].(map[string]interface{})

	policyRequest := client_iam.PolicyRequest{}

	if profile != "" {
		policyMap := getServicePolicyMapFromServiceType(d.Get("service").(string), profile)
		policyRequest = client_iam.PolicyRequest{
			Policy:  policyMap,
			Service: d.Get("service").(string),
		}
	} else { // custom policy
		policyMap := map[string]string{}
		for k, v := range custom {
			policyMap[k] = v.(string)
		}
		policyRequest = client_iam.PolicyRequest{
			Policy:  policyMap,
			Service: d.Get("service").(string),
		}
	}

	err := (*client_iam.Client)(c).AssignPolicyToUser(&policyRequest, d.Get("principal").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("principal").(string))

	return diags
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	c.HostURL = c.ServiceURLs["iam"]

	// delete policy
	policyRequest := client_iam.PolicyRequest{
		Service: d.Get("service").(string),
	}

	var diags diag.Diagnostics

	err := (*client_iam.Client)(c).RemovePolicyFromUser(&policyRequest, d.Get("principal").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getServicePolicyMapFromServiceType(serviceType string, profile string) map[string]string {
	policyMap := map[string]string{}
	switch serviceType {
	case "car-rental":
		policy := client_iam.CarrentalPolicy{}
		switch profile {
		case "reader":
			policy.DescribeCarOrders = "*"
			policy.ListCarOrders = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)
	case "ssl-monitoring":
		policy := client_iam.SslMonitoringPolicy{}
		switch profile {
		case "reader":
			policy.DescribeMonitors = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "azure-app-management":
		policy := client_iam.AzureAppManagementPolicy{}
		switch profile {
		case "reader":
			policy.ListCustomers = "*"
			policy.ListEnvironment = "*"
			policy.DescribeApplicationManagement = "*"
			policy.DescribeCustomer = "*"
			policy.DescribeEnvironment = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "azure-landing-zone":
		policy := client_iam.AzureLandingZonePolicy{}
		switch profile {
		case "reader":
			policy.ReadLandingZone = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "azureGovernance":
		policy := client_iam.AzureGovernancePolicy{}
		switch profile {
		case "reader":
			policy.DescribeCustomer = "*"
			policy.DescribeEnvironment = "*"
			policy.DescribeDeployment = "*"
			policy.DescribeTerracoreProject = "*"
			policy.DescribeTerracoreSchedule = "*"
			policy.ListCustomers = "*"
			policy.ListDataSync = "*"
			policy.ListEnvironment = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "cmsmgw":
		policy := client_iam.ManagedGatewayPolicy{}
		switch profile {
		case "reader":
			policy.ReadMgw = "*"
			policy.ReadConnection = "*"
			policy.ReadTranslation = "*"
			policy.ReadSpark = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "corc":
		policy := client_iam.ContainerOrchestrationPolicy{}
		switch profile {
		case "reader":
			policy.DescribeClusterUsers = "*"
			policy.DescribeClusters = "*"
			policy.ListClusterUsers = "*"
			policy.ListClusters = "*"
			policy.ListMyUsers = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "customer-db":
		policy := client_iam.CustomerDBPolicy{}
		switch profile {
		case "reader":
			policy.ReadAccess = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "domdns":
		policy := client_iam.DnsPolicy{}
		switch profile {
		case "reader":
			policy.DescribeRecords = "*"
			policy.DescribeZones = "*"
			policy.ListRecords = "*"
			policy.ListZones = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "drive":
		policy := client_iam.DrivePolicy{}
		switch profile {
		case "reader":
			policy.ReadAccess = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "event-bridge":
		policy := client_iam.EventBridgePolicy{}
		switch profile {
		case "reader":
			policy.ReadEventConfig = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "generic-services":
		policy := client_iam.GenericServicePolicy{}
		switch profile {
		case "reader":
			policy.ReadGeneric = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "iam":
		policy := client_iam.IamPolicy{}
		switch profile {
		case "reader":
			policy.AssumeRole = "*"
			policy.DescribePolicyDocument = "*"
			policy.DescribeRoles = "*"
			policy.DescribeServiceUsers = "*"
			policy.DescribeTrustRelations = "*"
			policy.DescribeTrustRelationsAdmin = "*"
			policy.DescribeUsers = "*"
			policy.ListPolicyDocument = "*"
			policy.ListRoles = "*"
			policy.ListServiceUsers = "*"
			policy.ListUsers = "*"
			policy.ReadSessions = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "managed-os":
		policy := client_iam.ManagedOSPolicy{}
		switch profile {
		case "reader":
			policy.DescribeVmwareInstances = "*"
			policy.ListVmwareInstances = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "metering-engine":
		policy := client_iam.MeteringEnginePolicy{}
		switch profile {
		case "reader":
			policy.DescribeMetering = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "mia-museum":
		policy := client_iam.MiaMuseumPolicy{}
		switch profile {
		case "reader":
			policy.ReadAccess = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "service-registry":
		policy := client_iam.ServiceRegistryPolicy{}
		switch profile {
		case "reader":
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "sms-notify":
		policy := client_iam.SmsNotificationsPolicy{}
		switch profile {
		case "reader":
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "t-rex":
		policy := client_iam.ResourceExplorerPolicy{}
		switch profile {
		case "reader":
			policy.Read = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	case "terracore":
		policy := client_iam.TerracorePolicy{}
		switch profile {
		case "reader":
			policy.DescribeProject = "*"
			policy.DescribeExecution = "*"
			policy.DescribeExecutionPlan = "*"
			policy.DescribeSchedule = "*"
			policy.ListProjects = "*"
			policy.ListSchedules = "*"
		case "full-access":
			setEveryValueInStructTo(&policy, "*")
		}
		policyMap = createMapFromStruct(policy)

	}
	return policyMap
}

func LowerInitial(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func setEveryValueInStructTo(s interface{}, value string) {
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		v.Field(i).SetString(value)
	}
}

func createMapFromStruct(s interface{}) map[string]string {
	v := reflect.ValueOf(s)
	policyMap := map[string]string{}
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldName = LowerInitial(fieldName)
		policyMap[fieldName] = v.Field(i).String()
	}
	return policyMap
}
