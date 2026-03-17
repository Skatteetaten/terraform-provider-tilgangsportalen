package tilgangsportalapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// ListEntraGroupsForRole lists the display names and Entra group (entitlement) UIDs
// of all Entra ID groups assigned to a role identified by RoleName
// See https://wiki.sits.no/spaces/OIM/pages/1380418936/15.+List+Azure+AD+Groups+For+Role+V1
func (client *Client) ListEntraGroupsForRole(RoleName string) (*EntraGroups, error) {
	var data EntraGroups
	log.Printf("Listing Entra Groups for Role %s ...", RoleName)
	// Construct the URL
	listEntraGroupsURL := fmt.Sprintf("/SKAT_RoleGovernance/ListAzureADGroupsForRoleV1?RoleName=%s", url.QueryEscape(RoleName))
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
