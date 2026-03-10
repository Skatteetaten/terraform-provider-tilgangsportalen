package tilgangsportalapi

import (
	"encoding/json"
	"log"
)

// ListEntraGroupsV1 lists the display name and Entra group (entitlement) ID
// of all Entra ID groups created/owned by the authenticated user
// See https://wiki.sits.no/spaces/OIM/pages/1380418941/14.+List+Azure+AD+Groups+V1
func (client *Client) ListEntraGroups() (*EntraGroups, error) {
	var data EntraGroups
	log.Println("Listing Entra groups...")
	// Construct the URL
	listEntraGroupsURL := "/SKAT_RoleGovernance/ListAzureADGroupsV1"
	// Perform the POST request
	response, err := client.GetRequest(listEntraGroupsURL)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the EntraGroups struct
	err = json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}

	log.Printf("Listing Entra groups successful. Found %d group(s).", len(data.EntraGroups))

	return &data, nil
}
