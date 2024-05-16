package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Updates the name of a role identified by its current name
// See https://wiki.sits.no/display/IDABAS/6.+Rename+Role
func (client *Client) RenameSystemRole(group RenameSystemRole) (*http.Response, error) {

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Renaming System Role from %s to %s...", group.OldName, group.NewName)

	// Construct the URL
	Rename_group_url := "/SKAT_RoleGovernance/RenameRole"

	// Perform the POST request
	response, err := client.PostRequest(Rename_group_url, nil, group_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Renaming System Role from %s to %s was successful", group.OldName, group.NewName)

	return response, nil
}
