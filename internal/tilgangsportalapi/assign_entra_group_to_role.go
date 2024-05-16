package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Assign an Entra group to a role
func (client *Client) AssignEntraGroupToRole(assignment EntraGroupRoleAssignment) (*http.Response, error) {

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(assignment)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Adding Entra Group %s to Role %s...", assignment.EntraGroup, assignment.RoleName)

	// Construct the URL
	create_group_url := "/SKAT_RoleGovernance/AssignAzureADEntitlementToRole"

	// Perform the POST request
	log.Printf("Performing POST request to url %s, with body %s",create_group_url,group_body)
	response, err := client.PostRequest(create_group_url, nil, group_body)
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
