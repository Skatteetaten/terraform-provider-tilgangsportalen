package tilgangsportalapi

import (
	"log"
)

// CheckIfGroupExists checks if an Entra group with display name groupName
// exists, and is owned by authenticated system user.
func (client *Client) CheckIfGroupExists(groupName string) (bool, error) {
	tempEntraGroups, err := client.ListEntraGroups()
	if err != nil {
		return false, err
	}

	for _, apiGroup := range tempEntraGroups.EntraGroups {
		if apiGroup.DisplayName == groupName {
			log.Printf("Group %s exists", groupName)
			return true, nil
		}
	}

	log.Printf("Group %s does not exist.", groupName)
	return false, nil
}
