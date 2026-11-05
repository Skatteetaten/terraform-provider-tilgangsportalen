package tilgangsportalapi

// CheckIfGroupExists checks if an Entra group with display name groupName
// exists, and is owned by authenticated system user (case-insensitive matching).
// Returns true if the group exists, the actual group name from the API, and any error encountered.
func (client *Client) CheckIfGroupExists(groupName string) (bool, string, error) {
	group, err := client.GetEntraGroup(groupName, false)

	if err != nil {
		// 605 (not found or no permissions) and 606 (not found) are answers, not failures: the group is not visible to us.
		if IsErrorDoesNotExistOrUnauthorized(err.Error()) || IsErrorDoesNotExist(err.Error()) {
			return false, "", nil
		}
		return false, "", err
	}

	return true, group.DisplayName, nil
}
