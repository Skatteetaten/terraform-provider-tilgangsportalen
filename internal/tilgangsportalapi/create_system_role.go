package tilgangsportalapi

import (
	"log"
	"net/http"
)

// CreateSystemRole creates a system role. The fields Name, ApprovalLevel,
// IsForITShop, ProductCategory are mandatory and need to be included in the
// SystemRole object. The name must be unique.
// See https://wiki.sits.no/display/IDABAS/2.1.+Create+Role
func (client *Client) CreateSystemRole(role SystemRole) (*http.Response, error) {

	var roleBody, err = CreateRequestBody(role)
	if err != nil {
		return nil, err
	}
	log.Printf("Creating role %s...", role.Name)

	// Construct the URL
	createRoleURL := "/SKAT_RoleGovernance/CreateRole"

	// Perform the POST request
	response, err := client.PostRequest(createRoleURL, roleBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Creation of System Role %s was successful.", role.Name)

	return response, nil
}
