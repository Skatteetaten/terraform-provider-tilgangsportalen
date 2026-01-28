package tilgangsportalapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

// GetSystemRoleId gets information about a specific role ID. 
// It gets the name, description, approval level, system role owner,
// system role security owner and product category for a specific
// system role.
// See https://wiki.sits.no/spaces/OIM/pages/1460147567/17.+Get+Role+ID
func (client *Client) GetSystemRoleID(objectID string) (*SystemRole, error) {
	var data SystemRole
	log.Printf("Fetching system role by ID %s", objectID)

	// Construct the URL, with query escape to handle special characters in role name
	getRoleURL := "/SKAT_RoleGovernance/GetRoleID?RoleUID=" + url.QueryEscape(objectID)
	
	// Perform the POST request
	response, err := client.GetRequest(getRoleURL)
	if err != nil {
		log.Printf("Role with ID \"%s\" was not found.", objectID)
		return nil, err
	}

	// Unmarshal the JSON data into the SystemRole struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		log.Printf("An error was thrown when unmarshaling response from tilgangsportalen.")
		return nil, err
	}

	// always set response for field L2Ident and L3Ident to lowercase as API returns uppercase
	data.L2Ident = strings.ToLower(data.L2Ident)
	data.L3Ident = strings.ToLower(data.L3Ident)

	log.Printf("Role with ID %s was found.", objectID)

	return &data, nil
}
