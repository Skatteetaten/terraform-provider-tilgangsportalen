package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Publishes an existing (created) system role identified by name to the 
// specified IT shop. This makes the role visible to users.
// See https://wiki.sits.no/display/IDABAS/5.+Publish+systemrole+to+IT-shop
func (client *Client) PublishSystemRole(role SystemRole) (*http.Response, error) {
	var role_body map[string]interface{}

	data := PublishSystemRole{
		Name:   role.Name,
		ITShop: role.ItShopName,
	}

	temp_body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &role_body)

	log.Printf("Publishing role %s...", role.Name)

	// Construct the URL
	publish_role_url := "/SKAT_RoleGovernance/PublishRoleToITShop"

	// Perform the POST request
	response, err := client.PostRequest(publish_role_url, nil, role_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Publishing of System Role %s was successful.", role.Name)

	return response, nil
}
