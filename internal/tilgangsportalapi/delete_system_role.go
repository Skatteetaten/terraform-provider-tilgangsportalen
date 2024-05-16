package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Deletes a system role identified by its name. If DeleteSystemRole.Force equals 1, the role will be
// deleted along with any role assignments it may have. If it is 0 the role will not be deleted if it
// has any role assignments.
// See https://wiki.sits.no/display/IDABAS/7.+Delete+Role
func (client *Client) DeleteSystemRole(role DeleteSystemRole) (*http.Response, error) {

	var role_body map[string]interface{}
	temp_body, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &role_body)

	log.Printf("Deleting role %s...", role.Name)

	// Construct the URL
	delete_role_url := "/SKAT_RoleGovernance/DeleteRole"

	// Perform the POST request
	response, err := client.PostRequest(delete_role_url, nil, role_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Deletion of System Role %s was successful.", role.Name)

	return response, nil
}
