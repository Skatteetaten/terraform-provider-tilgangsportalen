package tilgangsportalapi

import (
	"log"
	"net/http"
)

// AssignEntraGroupToRole creates a role assignment between a role and an
// Entra group
func (client *Client) AssignEntraGroupToRole(assignment EntraGroupRoleAssignment) (*http.Response, error) {

	var groupBody, err = CreateRequestBody(assignment)
	if err != nil {
		return nil, err
	}

	log.Printf("Adding Entra Group %s to Role %s...", assignment.EntraGroup, assignment.RoleName)

	// Construct the URL
	createGroupURL := "/SKAT_RoleGovernance/AssignAzureADEntitlementToRole"

	// Perform the POST request
	log.Printf("Performing POST request to url %s, with body %s", createGroupURL, groupBody)
	response, err := client.PostRequest(createGroupURL, groupBody)
	if err != nil {
		return nil, err
	}

	// Wait for assignment to be completed
	err = client.WaitForGroupRoleAssignmentStatus(assignment.EntraGroup, assignment.RoleName, true)
	if err != nil {
		return nil, err
	}

	log.Printf("Adding Entra Group %s to Role %s was successful!", assignment.EntraGroup, assignment.RoleName)

	return response, nil
}
