package tilgangsportalapi

import (
	"log"
	"net/http"
)

// UpdateSystemRole updates description, approval level, system role owner, 
// system role security owner, and/or product category for a specific (named)
// system role. To update role name see separate method.
// See https://wiki.sits.no/display/IDABAS/17.+Update+Role
func (client *Client) UpdateSystemRole(role SystemRoleChange) (*http.Response, error) {

	var roleBody, err = CreateRequestBody(role)
	if err != nil {
		return nil, err
	}

	log.Printf("Updating role %s...", role.RoleName)

	// Construct the URL
	updateRoleURL := "/SKAT_RoleGovernance/UpdateRole"

	// Perform the POST request
	response, err := client.PutRequest(updateRoleURL, roleBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Update of System Role %s was successful.", role.RoleName)

	return response, nil
}
