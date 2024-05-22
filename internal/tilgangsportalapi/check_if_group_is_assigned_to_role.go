package tilgangsportalapi

import (
	"log"
	"time"
)

// CheckIfGroupIsAssignedToRole - For a role with name roleName and a group with
// display name goupName, check if a role assignment exists between them.
func (client *Client) CheckIfGroupIsAssignedToRole(groupName string, roleName string) (bool, error) {
	tempEntraGroups, err := client.ListEntraGroupsForRole(roleName)
	if err != nil {
		return false, err
	}

	for _, apiGroup := range tempEntraGroups.EntraGroups {
		if apiGroup.DisplayName == groupName {
			log.Printf("Group %s is assigned to role %s", groupName, roleName)
			return true, nil
		}
	}

	log.Printf("Group %s is not assigned to role %s", groupName, roleName)
	return false, nil
}

// WaitForGroupRoleAssignmentStatus is used to verify that a role
// assignment has been either created (assignmentStatus=true) or removed
// (assignmentStatus=false). It uses ListEntraGroupsForRole in a while loop,
// when the assignment is either created or removed, the loop exits.
func (client *Client) WaitForGroupRoleAssignmentStatus(groupName string, roleName string, assignmentStatus bool) error {
	var action string

	if assignmentStatus {
		action = "assigned"
	} else {
		action = "removed"
	}

	log.Printf("Waiting for group %s to be %s to role %s...", groupName, action, roleName)

	for {
		var err error
		assignmentFound, err := client.CheckIfGroupIsAssignedToRole(groupName, roleName)
		if err != nil {
			return err
		}

		if assignmentStatus && assignmentFound {
			log.Printf("Group %s is assigned to role %s", groupName, roleName)
			break // Exit the loop if waiting for assignment and it is found
		} else if !assignmentStatus && !assignmentFound {
			log.Printf("Group %s is no longer assigned to role %s", groupName, roleName)
			break // Exit the loop if waiting for removal and it is removed
		}
		time.Sleep(5 * time.Second) // Wait before retrying

	}

	return nil
}
