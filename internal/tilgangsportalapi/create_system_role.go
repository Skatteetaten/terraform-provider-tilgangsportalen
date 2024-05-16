package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/http"
)

// Creates a system role. The fields Name, ApprovalLevel, IsForITShop, ProductCategory are 
// mandatory and need to be included in the SystemRole object. The name must be unique.
// See https://wiki.sits.no/display/IDABAS/2.1.+Create+Role
func (client *Client) CreateSystemRole(role SystemRole) (*http.Response, error) {

	var role_body map[string]interface{}
	temp_body, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(temp_body, &role_body)

	log.Printf("Creating role %s...", role.Name)

	// Construct the URL
	create_role_url := "/SKAT_RoleGovernance/CreateRole"

	// Perform the POST request
	response, err := client.PostRequest(create_role_url, nil, role_body)
	if err != nil {
		return nil, err
	}

	log.Printf("Creation of System Role %s was successful.", role.Name)

	return response, nil
}
