package dynamiccloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dynamiccloud "github.com/cancom/terraform-provider-cancom/client/services/dynamic-cloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVpcProjectUsers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcProjectUsersCreate,
		ReadContext:   resourceVpcProjectUsersRead,
		UpdateContext: resourceVpcProjectUsersUpdate,
		DeleteContext: resourceVpcProjectUsersDelete,
		CustomizeDiff: resourceVpcProjectUsersCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the tenant this VPC Project belongs to.",
			},
			"users": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The list of users with access to the VPC Project. The list may only contains CRNs of human iam users.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					// the tf plugin sdk only allows to validate each value and not the entire list, so CustomizeDiff is used to validate the list is not empty
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(CrnRegex, "One of the users is not a valid CANCOM Resource Numbers (CRNs).")),
				},
			},
			"vpc_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id (short uuid) of the VPC Project.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceVpcProjectUsersCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	usersSet := d.Get("users").(*schema.Set)
	if usersSet.Len() == 0 {
		return fmt.Errorf("invalid value for users (the set of users must include at least one CRN)")
	}
	return nil
}

func resourceVpcProjectUsersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "Creating VPC Project users")

	return resourceVpcProjectUsersUpdate(ctx, d, meta)
}

func resourceVpcProjectUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Reading VPC Project users")
	var diags diag.Diagnostics
	vpcProjectShortid, found := strings.CutSuffix(d.Id(), "_users")
	if !found {
		return diag.Errorf("error parsing VPC Project id")
	}

	// GetVpcProject returns nil if the VPC Project is NotFound
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		d.SetId("")
		return nil
	}
	fullProjectName := fmt.Sprintf("%s-%s", resp.Metadata.Tenant, resp.Metadata.Name)
	humanUsers, err := GetHumanUsers(fullProjectName, resp.Spec.ProjectUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(humanUsers) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("vpc_project_id", resp.Metadata.Shortid)
	d.Set("tenant", resp.Metadata.Tenant)

	err = d.Set("users", usersToSet(humanUsers))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceVpcProjectUsersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Updating VPC Project users")
	vpcProjectShortid := d.Get("vpc_project_id").(string)
	tflog.Info(ctx, "Parsed vpc_project_id", map[string]interface{}{"vpcProjectShortid": vpcProjectShortid})
	d.SetId(fmt.Sprintf("%s_users", vpcProjectShortid))
	users := setToUsers(d.Get("users").(*schema.Set))

	// get VPC Project to get any serviceUsers created in the project
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		d.SetId("")
		return diag.Errorf("error updating VPC Project users. VPC Project NotFound")
	}
	fullProjectName := fmt.Sprintf("%s-%s", resp.Metadata.Tenant, resp.Metadata.Name)
	serviceUsers, err := GetServiceUsers(fullProjectName, resp.Spec.ProjectUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Concat svc users", map[string]interface{}{"usersForBody": append(users, serviceUsers...)})
	_, err = (*client_dynamiccloud.Client)(c).UpdateVpcProjectUsers(vpcProjectShortid, append(users, serviceUsers...))
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *resource.RetryError {
		return WaitProjectReady(ctx, c, vpcProjectShortid)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcProjectUsersRead(ctx, d, meta)
}

func resourceVpcProjectUsersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := meta.(*client.CcpClient).GetService("dynamic-cloud")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Deleting VPC Project users")
	var diags diag.Diagnostics
	vpcProjectShortid, found := strings.CutSuffix(d.Id(), "_users")
	if !found {
		return diag.Errorf("error parsing VPC Project id")
	}

	// get VPC Project to add any serviceUsers already created in the project
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(vpcProjectShortid)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		d.SetId("")
		return nil
	}
	fullProjectName := fmt.Sprintf("%s-%s", resp.Metadata.Tenant, resp.Metadata.Name)
	serviceUsers, err := GetServiceUsers(fullProjectName, resp.Spec.ProjectUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = (*client_dynamiccloud.Client)(c).UpdateVpcProjectUsers(vpcProjectShortid, serviceUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete)-time.Minute, func() *resource.RetryError {
		return WaitProjectReady(ctx, c, vpcProjectShortid)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}

func usersToSet(userList []string) *schema.Set {
	interfaceSlice := make([]interface{}, len(userList))
	for i, user := range userList {
		interfaceSlice[i] = user
	}
	users := schema.NewSet(schema.HashString, interfaceSlice)
	return users
}

func setToUsers(users *schema.Set) []string {
	userList := make([]string, 0, len(users.List()))
	for _, user := range users.List() {
		userList = append(userList, user.(string))
	}
	if userList == nil {
		userList = []string{}
	}
	return userList
}
