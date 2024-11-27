package dynamiccloud

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var CrnIamUserRegex = regexp.MustCompile("^crn:[A-Za-z0-9]+:[A-Za-z0-9]*:iam:user:[A-Za-z0-9-.]+@[A-Za-z0-9-.]+$")

func IsIamUserCrn(crn string) bool {
	return CrnIamUserRegex.MatchString(crn)
}

func isVPCProjectServiceUserName(userName, tenant, projectName string) bool {
	fullProjectName := fmt.Sprintf("%s-%s", tenant, projectName)
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
		splitUserName[1] == fullProjectName &&
		numberSuffix > 0 &&
		numberSuffix < 10000)
}

func getHumanUsers(tenant, projectName string, users []string) ([]string, error) {
	humanUsers := []string{}
	for _, user := range users {
		if IsIamUserCrn(user) {
			humanUsers = append(humanUsers, user)
		} else if !isVPCProjectServiceUserName(user, tenant, projectName) {
			return nil, fmt.Errorf("error parsing users - user %s is not a service user or CRN", user)
		}
	}
	return humanUsers, nil
}

func getServiceUsers(tenant, projectName string, users []string) ([]string, error) {
	serviceUsers := []string{}
	for _, user := range users {
		if isVPCProjectServiceUserName(user, tenant, projectName) {
			serviceUsers = append(serviceUsers, user)
		} else if !IsIamUserCrn(user) {
			return nil, fmt.Errorf("error parsing users - user %s is not a service user or CRN", user)
		}
	}
	return serviceUsers, nil
}

func usersToSet(userList []string) *schema.Set {
	interfaceSlice := make([]interface{}, len(userList))
	for i, user := range userList {
		interfaceSlice[i] = user
	}
	users := schema.NewSet(schema.HashString, interfaceSlice)
	return users
}

func setToUsers(users interface{}) []string {
	if users == nil {
		return []string{}
	}
	usersSet := users.(*schema.Set)
	userList := make([]string, 0, len(usersSet.List()))
	for _, user := range usersSet.List() {
		userList = append(userList, user.(string))
	}
	if userList == nil {
		userList = []string{}
	}
	return userList
}
