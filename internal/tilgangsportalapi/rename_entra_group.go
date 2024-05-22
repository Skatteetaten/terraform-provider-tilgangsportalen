package tilgangsportalapi

import (
	"log"
	"net/http"
)

// RenameEntraGroup updates the name of an Entra ID group identified by its
// current name
// See https://wiki.sits.no/display/IDABAS/9.+Rename+Azure+AD+Group
func (client *Client) RenameEntraGroup(group RenameEntraGroup) (*http.Response, error) {

	var groupBody, err = CreateRequestBody(group)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming Entra group %s to %s...", group.OldName, group.NewName)

	// Construct the URL
	renameGroupURL := "/SKAT_RoleGovernance/RenameAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(renameGroupURL, groupBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming Entra group %s to %s was successful.", group.OldName, group.NewName)

	return response, nil
}
