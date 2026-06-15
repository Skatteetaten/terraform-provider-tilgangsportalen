package tilgangsportalapi

import (
	"context"
	"fmt"
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

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

		select {
		case <-time.After(5 * time.Second): // Retry after 5 seconds
		case <-ctx.Done(): // Return error if timeout is reached
			// Wrap ctx.Err() so callers can detect timeout/cancel, while still getting a readable message with lookup context.
			timeoutErr := fmt.Errorf("timeout reached while waiting for group %s to be %s to role %s: %w", groupName, action, roleName, ctx.Err())
			return timeoutErr
		}
	}

	return nil
}
