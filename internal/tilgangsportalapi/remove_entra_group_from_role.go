package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Removes an Entra ID group identified by name, from a system role identified
// by name.  
// See https://wiki.sits.no/display/IDABAS/12.+RemoveAzureADEntitlementFromRole
func (client *Client) RemoveEntraGroupFromRole(group EntraGroupRoleAssignment) (*http.Response, error) {

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Removing Entra Group %s from Role %s...", group.EntraGroup, group.RoleName)

	// Construct the URL
	create_group_url := "/SKAT_RoleGovernance/RemoveAzureADEntitlementFromRole"

	// Perform the POST request
	response, err := client.PostRequest(create_group_url, nil, group_body)
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
