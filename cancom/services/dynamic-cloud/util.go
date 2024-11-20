package dynamiccloud

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cancom/terraform-provider-cancom/client"
	client_dynamiccloud "github.com/cancom/terraform-provider-cancom/client/services/dynamic-cloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var CrnRegex = regexp.MustCompile("^crn:[A-Za-z0-9]+:[A-Za-z0-9]*:[A-Za-z0-9-]+:[A-Za-z0-9-]+:[A-Za-z0-9-@.]+$")

func GetHumanUsers(projectName string, users []string) ([]string, error) {
	humanUsers := []string{}
	for _, user := range users {
		if IsUserCrn(user) {
			humanUsers = append(humanUsers, user)
		} else if !IsServiceUserName(user, projectName) {
			return nil, fmt.Errorf("error parsing users - user %s is not a service user or CRN", user)
		}
	}
	return humanUsers, nil
}

func GetServiceUsers(projectName string, users []string) ([]string, error) {
	serviceUsers := []string{}
	for _, user := range users {
		if IsServiceUserName(user, projectName) {
			serviceUsers = append(serviceUsers, user)
		} else if !IsUserCrn(user) {
			return nil, fmt.Errorf("error parsing users - user %s is not a service user or CRN", user)
		}
	}
	return serviceUsers, nil
}

func IsServiceUserName(userName string, projectName string) bool {
	if strings.Contains(userName, "@") {
		return false
	}
	splitUserName := strings.SplitN(userName, "_", 3)
	if len(splitUserName) != 3 {
		return false
	}
	numberSuffix, err := strconv.Atoi(splitUserName[2])
	if err != nil {
		return false
	}
	return (splitUserName[0] == "svc" &&
		splitUserName[1] == projectName &&
		numberSuffix > 0 &&
		numberSuffix < 10000)
}

func IsUserCrn(crn string) bool {
	if !CrnRegex.MatchString(crn) {
		return false
	}
	splitCrn := strings.SplitN(crn, ":", 6)
	return (splitCrn[3] == "iam" &&
		splitCrn[4] == "user" &&
		strings.Contains(splitCrn[5], "@"))
}

func waitProject(ctx context.Context, c *client.Client, project_id string) (string, *resource.RetryError) {
	// GetVpcProject returns nil if the VPC Project is NotFound
	resp, err := (*client_dynamiccloud.Client)(c).GetVpcProject(project_id)
	if err != nil {
		return "", resource.NonRetryableError(fmt.Errorf("error describing VPC Project: %s", err))
	}
	if resp == nil {
		return "NotFound", resource.NonRetryableError(fmt.Errorf("error describing VPC Project. VPC Project 'NotFound'"))
	}

	switch resp.Status.Phase {
	case "Error":
		return "", resource.NonRetryableError(fmt.Errorf("error describing VPC Project. status.phase is Error"))
	case "Ready":
		return "Ready", nil
	default:
		tflog.Info(ctx, fmt.Sprintf("Waiting for VPC Project to finish %s", resp.Status.Phase))
		return "", resource.RetryableError(fmt.Errorf("VPC Project is still transitioning with phase '%s'", resp.Status.Phase))
	}
}

func WaitProjectReady(ctx context.Context, c *client.Client, project_id string) *resource.RetryError {
	phase, err := waitProject(ctx, c, project_id)
	if err != nil {
		return err
	}
	if phase != "Ready" {
		return resource.NonRetryableError(fmt.Errorf("expected VPC Project status phase to be 'Ready' but was in phase '%s'", phase))
	}
	return nil
}

func WaitProjectDeleted(ctx context.Context, c *client.Client, project_id string) *resource.RetryError {
	phase, err := waitProject(ctx, c, project_id)
	if phase == "NotFound" {
		return nil
	}
	if err != nil {
		return err
	}
	return resource.NonRetryableError(fmt.Errorf("expected VPC Project status phase to be 'NotFound' but was in phase '%s'", phase))
}
