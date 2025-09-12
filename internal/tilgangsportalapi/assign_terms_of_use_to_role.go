package tilgangsportalapi

import (
	"log"
	"net/http"
)

// AssignTermsOfUseToRole assigns a terms of use object to a role
// See https://wiki.sits.no/spaces/OIM/pages/1333270042/24.+Assign+Terms+of+Use+to+Role
func (client *Client) AssignTermsOfUseToRole(assignment TermsOfUseAssignment) (*http.Response, error) {

	var termsOfUseAssignmentBody, err = CreateRequestBody(assignment)
	if err != nil {
		return nil, err
	}

	log.Printf("Adding Terms of Use %s to Role %s...", assignment.TermsOfUse, assignment.RoleName)

	// Construct the URL
	assignTermsOfUseURL := "/SKAT_RoleGovernance/AssignTermsOfUseToRole"

	// Perform the POST request
	log.Printf("Performing POST request to url %s, with body %s", assignTermsOfUseURL, termsOfUseAssignmentBody)
	response, err := client.PostRequest(assignTermsOfUseURL, termsOfUseAssignmentBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Adding Terms of Use %s to Role %s was successful!", assignment.TermsOfUse, assignment.RoleName)

	return response, nil
}
