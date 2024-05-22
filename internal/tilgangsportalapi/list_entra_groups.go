package tilgangsportalapi

import (
	"encoding/json"
	"log"
)

// ListEntraGroups lists the display names of all Entra ID groups created/owned
// by the authenticated user
// See https://wiki.sits.no/display/IDABAS/14.+ListAzureADGroups
func (client *Client) ListEntraGroups() (*EntraGroups, error) {
	var data EntraGroups
	log.Println("Listing Entra groups...")
	// Construct the URL
	listEntraGroupsURL := "/SKAT_RoleGovernance/ListAzureADGroups"
	// Perform the POST request
	response, err := client.GetRequest(listEntraGroupsURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the Roles struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	log.Printf("Listing Entra groups successful. Found %d group(s).", len(data.EntraGroups))

	return &data, nil
}
