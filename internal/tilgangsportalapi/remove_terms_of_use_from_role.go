package tilgangsportalapi

import (
	"log"
	"net/http"
)

// RemoveTermsOfUseFromRole removes a terms of use object from a role
// See https://wiki.sits.no/spaces/OIM/pages/1333270058/25.+Remove+Terms+of+Use+from+Role
func (client *Client) RemoveTermsOfUseFromRole(role TermsOfUseRemoval) (*http.Response, error) {

	var termsOfUseRemovalBody, err = CreateRequestBody(role)
	if err != nil {
		return nil, err
	}

	log.Printf("Removing Terms of Use from Role %s...", role.RoleName)

	// Construct the URL
	removeTermsOfUseURL := "/SKAT_RoleGovernance/RemoveTermsOfUseFromRole"

	// Perform the POST request
	response, err := client.PostRequest(removeTermsOfUseURL, termsOfUseRemovalBody)
	if err != nil {
		return nil, err
	}
	// wait for assignment to be completed
	err = client.WaitForTermsOfUseRemoval(role.RoleName)
	if err != nil {
		return nil, err
	}

	log.Printf("Removing Terms of Use from Role %s was successful.", role.RoleName)

	return response, nil
}
