package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Updates description, approval level, system role owner, system role security owner, 
// and/or product category for a specific (named) system role. To update role name see 
// separate method.
// See https://wiki.sits.no/display/IDABAS/17.+Update+Role
func (client *Client) UpdateSystemRole(role SystemRoleChange) (*http.Response, error) {
	var role_body map[string]interface{}
	temp_body, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &role_body)

	log.Printf("Updating role %s...", role.RoleName)

	// Construct the URL
	update_role_url := "/SKAT_RoleGovernance/UpdateRole"

	// Perform the POST request
	response, err := client.PutRequest(update_role_url, nil, role_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Update of System Role %s was successful.", role.RoleName)

	return response, nil
}