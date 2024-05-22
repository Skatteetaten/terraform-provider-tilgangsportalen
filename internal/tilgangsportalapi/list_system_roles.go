package tilgangsportalapi

import (
	"encoding/json"
	"log"
)

// ListSystemRoles lists the display names of all roles created by the authenticated user
// See https://wiki.sits.no/display/IDABAS/13.+ListRoles
func (client *Client) ListSystemRoles() (*Roles, error) {
	var data Roles
	log.Println("Listing roles...")
	// Construct the URL
	listRolesURL := "/SKAT_RoleGovernance/ListRoles"

	// Perform the POST request
	response, err := client.GetRequest(listRolesURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the Roles struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	log.Printf("Listing roles successful. Found %d role(s).",len(data.Roles))

	return &data, nil
}
