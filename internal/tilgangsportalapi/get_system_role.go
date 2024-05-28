package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/url"
)

// GetSystemRole gets information about a specific named role / check if role
// exists. Gets the name, description, approval level, system role owner,
// system role security owner, product category and IsForITShop for a specific
// (named) system role.
// See https://wiki.sits.no/display/IDABAS/19.+Get+Role
func (client *Client) GetSystemRole(roleName string) (*SystemRole, error) {
	var data SystemRole
	log.Printf("Fetching system role %s", roleName)
	// Construct the URL, with query escape to handle special characters in role name
	getRoleURL := "/SKAT_RoleGovernance/GetRole?roleName=" + url.QueryEscape(roleName)
	// Perform the POST request
	response, err := client.GetRequest(getRoleURL)
	if err != nil {
		log.Printf("Role with name \"%s\" was not found.", roleName)
		return nil, err
	}

	// Unmarshal the JSON data into the SystemRole struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		log.Printf("An error was thrown when unmarshaling response from tilgangsportalen.")
		return nil, err
	}

	log.Printf("Role with name %s was found.", roleName)

	return &data, nil
}
