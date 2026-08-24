package staff

import (
	"fmt"
	"strings"
)

type ActionPolicy struct {
	Action string
	Roles  []string
}

var defaultPolicies = []ActionPolicy{
	{Action: "catalog", Roles: []string{RoleManager, RoleRentalDesk, RoleWarehouse}},
	{Action: "rental", Roles: []string{RoleManager, RoleRentalDesk}},
	{Action: "maintenance", Roles: []string{RoleManager, RoleTechnician}},
}

func CheckRole(role, action string) error {
	role = strings.TrimSpace(role)
	for _, policy := range defaultPolicies {
		if policy.Action != action {
			continue
		}
		for _, allowed := range policy.Roles {
			if role == allowed {
				return nil
			}
		}
		return fmt.Errorf("role %s cannot perform %s", role, action)
	}
	return fmt.Errorf("unknown action %s", action)
}

func Policies() []ActionPolicy {
	result := make([]ActionPolicy, len(defaultPolicies))
	copy(result, defaultPolicies)
	return result
}
