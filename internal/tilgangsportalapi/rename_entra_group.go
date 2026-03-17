package tilgangsportalapi

import (
	"encoding/json"
	"log"
)

// RenameEntraGroup updates the name of an Entra ID group identified by its
// current name
// See https://wiki.sits.no/spaces/OIM/pages/1044589610/9.+Rename+Azure+AD+Group
func (client *Client) RenameEntraGroup(group RenameEntraGroup) (*SuccessfulResponse, error) {

	var renameGroupBody, err = CreateRequestBody(group)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming Entra group %s to %s...", group.OldName, group.NewName)

	// Construct the URL
	renameGroupURL := "/SKAT_RoleGovernance/RenameAzureADGroup"

	// Perform the POST request
	urlRequestStr := client.baseURL + renameGroupURL
	_, responseBody, err := BuildRequest("POST", urlRequestStr, renameGroupBody, client.headers, client.cookies, defaultTimeout)
	if err != nil {
		return nil, err
	}

	var response SuccessfulResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming Entra group %s to %s was successful.", group.OldName, group.NewName)

	return &response, nil
}
