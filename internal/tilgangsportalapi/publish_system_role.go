package tilgangsportalapi

import (
	"log"
	"net/http"
)

// PublishSystemRole publishes an existing (created) system role identified by
// name to the specified IT shop. This makes the role visible to users.
// See https://wiki.sits.no/display/IDABAS/5.+Publish+systemrole+to+IT-shop
func (client *Client) PublishSystemRole(role SystemRole) (*http.Response, error) {

	publishRole := PublishSystemRole{
		Name:   role.Name,
		ITShop: role.ItShopName,
	}

	var roleBody, err = CreateRequestBody(publishRole)
	if err != nil {
		return nil, err
	}

	log.Printf("Publishing role %s...", role.Name)

	// Construct the URL
	publishRoleURL := "/SKAT_RoleGovernance/PublishRoleToITShop"

	// Perform the POST request
	response, err := client.PostRequest(publishRoleURL, roleBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Publishing of System Role %s was successful.", role.Name)

	return response, nil
}
