package tilgangsportalapi

import (
	"log"
	"strings"
)

// CheckIfGroupExists checks if an Entra group with display name groupName
// exists, and is owned by authenticated system user (case-insensitive matching).
// Returns true if the group exists, the actual group name from the API, and any error.
func (client *Client) CheckIfGroupExists(groupName string) (bool, string, error) {
	tempEntraGroups, err := client.ListEntraGroups()
	if err != nil {
		return false, "", err
	}

	for _, apiGroup := range tempEntraGroups.EntraGroups {
		if strings.EqualFold(apiGroup.DisplayName, groupName) {
			log.Printf("Group %s exists", groupName)
			return true, apiGroup.DisplayName, nil
		}
	}

	log.Printf("Group %s does not exist.", groupName)
	return false, "", nil
}
