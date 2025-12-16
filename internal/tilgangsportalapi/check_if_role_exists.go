package tilgangsportalapi

import "strings"

// CheckIfRoleExists calls GetSystemRole to check if a role with the name
// roleName exists (case-insensitive matching). Returns true if the role exists,
// the actual role name from the API, and any error encountered.
func (client *Client) CheckIfRoleExists(roleName string) (bool, string, error) {
	roles, err := client.ListSystemRoles()

	if err != nil {
		return false, "", err
	}

	for _, role := range roles.Roles {
		if strings.EqualFold(role.DisplayName, roleName) {
			return true, role.DisplayName, nil
		}
	}

	return false, "", nil
}
