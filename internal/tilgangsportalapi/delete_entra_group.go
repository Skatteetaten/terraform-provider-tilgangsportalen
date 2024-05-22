package tilgangsportalapi

import (
	"log"
	"net/http"
)

// DeleteEntraGroup deletes a system role identified by its name. If
// DeleteEntraGroup.Force equals 1, the group will be deleted along with any
// account assignments it may have. If it is 0 the group will not be deleted if
// it has any account assignments.
// See https://wiki.sits.no/display/IDABAS/8.+Delete+Azure+AD+Group
func (client *Client) DeleteEntraGroup(group DeleteEntraGroup) (*http.Response, error) {

	var groupBody, err = CreateRequestBody(group)
	if err != nil {
		return nil, err
	}

	log.Printf("Deleting Entra group %s...", group.Name)

	// Construct the URL
	deleteGroupURL := "/SKAT_RoleGovernance/DeleteAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(deleteGroupURL, groupBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Deletion of Entra group %s was successful.", group.Name)

	return response, nil
}
