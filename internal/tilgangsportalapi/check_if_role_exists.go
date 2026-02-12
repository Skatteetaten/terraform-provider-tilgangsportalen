package tilgangsportalapi

// CheckIfRoleExists checks if a role with the name roleName exists, and is owned
// by authenticated system user (case-insensitive matching).
// Returns true if the role exists, the actual role name from the API, and any error encountered.
func (client *Client) CheckIfRoleExists(roleName string) (bool, string, error) {
	role, err := client.GetSystemRole(roleName)

	if err != nil {
		// Role does not exist
		return false, "", err
	}

	return true, role.Name, nil
}
