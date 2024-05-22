package tilgangsportalapi

import (
	"log"
	"net/http"
)

// RemoveEntraGroupFromRole removes an Entra ID group identified by name, from
// a system role identified by name.
// See https://wiki.sits.no/display/IDABAS/12.+RemoveAzureADEntitlementFromRole
func (client *Client) RemoveEntraGroupFromRole(group EntraGroupRoleAssignment) (*http.Response, error) {

	var groupBody, err = CreateRequestBody(group)
	if err != nil {
		return nil, err
	}

	log.Printf("Removing Entra Group %s from Role %s...", group.EntraGroup, group.RoleName)

	// Construct the URL
	createGroupURL := "/SKAT_RoleGovernance/RemoveAzureADEntitlementFromRole"

	// Perform the POST request
	response, err := client.PostRequest(createGroupURL, groupBody)
	if err != nil {
		return nil, err
	}
	// wait for assignment to be completed
	err = client.WaitForGroupRoleAssignmentStatus(group.EntraGroup, group.RoleName, false)
	if err != nil {
		return nil, err
	}

	log.Printf("Removing Entra Group %s from Role %s was successful.", group.EntraGroup, group.RoleName)

	return response, nil
}
