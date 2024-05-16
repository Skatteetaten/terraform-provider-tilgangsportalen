package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Updates the name of an Entra ID group identified by its current name
// See https://wiki.sits.no/display/IDABAS/9.+Rename+Azure+AD+Group
func (client *Client) RenameEntraGroup(group RenameEntraGroup) (*http.Response, error) {

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Renaming Entra group %s to %s...", group.OldName, group.NewName)

	// Construct the URL
	Rename_group_url := "/SKAT_RoleGovernance/RenameAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(Rename_group_url, nil, group_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming Entra group %s to %s was successful.", group.OldName, group.NewName)

	return response, nil
}
