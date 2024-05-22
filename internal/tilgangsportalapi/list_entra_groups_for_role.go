package tilgangsportalapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// ListEntraGroupsForRole lists the display names of all Entra ID groups
// assigned to a role identified by RoleName
// See https://wiki.sits.no/display/IDABAS/15.+ListAzureADGroupsForRole
func (client *Client) ListEntraGroupsForRole(RoleName string) (*EntraGroups, error) {
	var data EntraGroups
	log.Printf("Listing Entra Groups for Role %s ...", RoleName)
	// Construct the URL
	listEntraGroupsURL := fmt.Sprintf("/SKAT_RoleGovernance/ListAzureADGroupsForRole?RoleName=%s", url.QueryEscape(RoleName))
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

	log.Printf("Listing Entra Groups for Role %s successful.", RoleName)

	return &data, nil
}
