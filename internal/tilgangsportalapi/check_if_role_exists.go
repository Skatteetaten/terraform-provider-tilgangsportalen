package tilgangsportalapi

// CheckIfRoleExists calls GetSystemRole to check if a role with the name
// roleName exists
func (client *Client) CheckIfRoleExists(roleName string) (bool, error) {
	roles, err := client.ListSystemRoles()

	if err != nil {
		return false, nil
	}

	for _, role := range roles.Roles {
		if role.DisplayName == roleName {
			return true, nil

		}
	}

	return false, nil
}
