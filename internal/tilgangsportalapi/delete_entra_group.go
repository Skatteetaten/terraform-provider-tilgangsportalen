package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Deletes a system role identified by its name. If DeleteEntraGroup.Force equals 1, the group will be 
// deleted along with any account assignments it may have. If it is 0 the group will not be deleted if it 
// has any account assignments.
// See https://wiki.sits.no/display/IDABAS/8.+Delete+Azure+AD+Group
func (client *Client) DeleteEntraGroup(group DeleteEntraGroup) (*http.Response, error) {

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Deleting Entra group %s...", group.Name)

	// Construct the URL
	delete_group_url := "/SKAT_RoleGovernance/DeleteAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(delete_group_url, nil, group_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Deletion of Entra group %s was successful.", group.Name)

	return response, nil
}
