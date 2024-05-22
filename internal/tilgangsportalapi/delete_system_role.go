package tilgangsportalapi

import (
	"log"
	"net/http"
)

// DeleteSystemRole deletes a system role identified by its name. If
// DeleteSystemRole.Force equals 1, the role will be deleted along with any
// role assignments it may have. If it is 0 the role will not be deleted if it
// has any role assignments.
// See https://wiki.sits.no/display/IDABAS/7.+Delete+Role
func (client *Client) DeleteSystemRole(role DeleteSystemRole) (*http.Response, error) {

	var roleBody, err = CreateRequestBody(role)
	if err != nil {
		return nil, err
	}

	log.Printf("Deleting role %s...", role.Name)

	// Construct the URL
	deleteRoleURL := "/SKAT_RoleGovernance/DeleteRole"

	// Perform the POST request
	response, err := client.PostRequest(deleteRoleURL, roleBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Deletion of System Role %s was successful.", role.Name)

	return response, nil
}
