package tilgangsportalapi

import (
	"log"
	"net/http"
)

// RenameSystemRole updates the name of a role identified by its current name
// See https://wiki.sits.no/display/IDABAS/6.+Rename+Role
func (client *Client) RenameSystemRole(role RenameSystemRole) (*http.Response, error) {

	var roleBody, err = CreateRequestBody(role)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming System Role from %s to %s...", role.OldName, role.NewName)

	// Construct the URL
	renameRoleURL := "/SKAT_RoleGovernance/RenameRole"

	// Perform the POST request
	response, err := client.PostRequest(renameRoleURL, roleBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming System Role from %s to %s was successful", role.OldName, role.NewName)

	return response, nil
}
