package tilgangsportalapi

import (
	"encoding/json"
	"log"
)

// CreateEntraGroup creates an Entra ID group. The fields DisplayName
// and Tenant are mandatory and need to be included in the EntraGroup object.
// The Displayname must be unique.
// See https://wiki.sits.no/display/IDABAS/2.+Create+systemrole
func (client *Client) CreateEntraGroup(group EntraGroup) (*SuccessfulResponse, error) {

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
	urlRequestStr := client.baseURL + createGroupURL
	_, responseBody, err := BuildRequest("POST", urlRequestStr, groupBody, client.headers, client.cookies, defaultTimeout)
	if err != nil {
		return nil, err
	}

	var response SuccessfulResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	log.Printf("Creation of Entra group %s was successful.", group.DisplayName)

	return &response, nil
}
