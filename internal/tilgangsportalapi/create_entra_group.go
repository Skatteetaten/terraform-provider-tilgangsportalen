package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Creates an Entra ID group. The fields DisplayName, Alias, and Tenant are mandatory and 
// need to be included in the EntraGroup object. The Displayname must be unique.
// See https://wiki.sits.no/display/IDABAS/2.+Create+systemrole
func (client *Client) CreateEntraGroup(group EntraGroup) (*http.Response, error) {

	// add static values to the group object: Inheritance_level and Tenant
	group.Tenant = "Skatteetaten"

	var group_body map[string]interface{}
	temp_body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &group_body)

	log.Printf("Creating Entra group %s...", group.DisplayName)

	// Construct the URL
	create_group_url := "/SKAT_RoleGovernance/CreateAzureADGroup"

	// Perform the POST request
	response, err := client.PostRequest(create_group_url, nil, group_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Creation of Entra group %s was successful.", group.DisplayName)

	return response, nil
}
