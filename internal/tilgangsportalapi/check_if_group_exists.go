package tilgangsportalapi

// CheckIfGroupExists checks if an Entra group with display name groupName
// exists, and is owned by authenticated system user (case-insensitive matching).
// Returns true if the group exists, the actual group name from the API, and any error encountered.
func (client *Client) CheckIfGroupExists(groupName string) (bool, string, error) {
	group, err := client.GetEntraGroup(groupName, false)

	if err != nil {
		// Group does not exist
		return false, "", err
	}

	return true, group.DisplayName, nil
}
