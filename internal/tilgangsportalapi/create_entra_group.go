package tilgangsportalapi

import (
	"log"
	"net/http"
)

// CreateEntraGroup creates an Entra ID group. The fields DisplayName
// and Tenant are mandatory and need to be included in the EntraGroup object.
// The Displayname must be unique.
// See https://wiki.sits.no/display/IDABAS/2.+Create+systemrole
func (client *Client) CreateEntraGroup(group EntraGroup) (*http.Response, error) {

	// add static values to the group object: Tenant
	group.Tenant = "Skatteetaten"

	var groupBody, err = CreateRequestBody(group)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating Entra group %s...", group.DisplayName)

	// Construct the URL
	createGroupURL := "/SKAT_RoleGovernance/CreateAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(createGroupURL, groupBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Creation of Entra group %s was successful.", group.DisplayName)

	return response, nil
}
